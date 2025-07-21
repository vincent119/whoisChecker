package internal

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/spf13/cobra"
)

var (
    domain     string
    domainFile string
    showASN    bool
    useIPv4    bool
    useIPv6    bool
    outputJSON bool
    outputCSV  bool
    workers    int
    verbose    bool
    retry      int
    delay      int
)

// WhoisResult represents the WHOIS query result
type WhoisResult struct {
    Domain     string `json:"domain" csv:"domain"`
    Registrar  string `json:"registrar" csv:"registrar"`
    Expires    string `json:"expires" csv:"expires"`
    DaysLeft   string `json:"days_left" csv:"days_left"`
    Status     string `json:"status" csv:"status"`
    Error      string `json:"error,omitempty" csv:"error"`
    ASN        string `json:"asn,omitempty" csv:"asn"`
    IPv4       string `json:"ipv4,omitempty" csv:"ipv4"`
    IPv6       string `json:"ipv6,omitempty" csv:"ipv6"`
    Country    string `json:"country,omitempty" csv:"country"`
}

var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("whois called")
	},
}

func init() {
    whoisCmd := &cobra.Command{
        Use:   "whois",
        Short: "Query WHOIS information",
        Run:   runWhois,
    }

    whoisCmd.Flags().StringVarP(&domain, "domain", "d", "", "Query single domain")
    whoisCmd.Flags().StringVar(&domainFile, "domainfile", "", "Read domain list from file")
    whoisCmd.Flags().BoolVarP(&showASN, "asn", "A", false, "Show ASN information")
    whoisCmd.Flags().BoolVar(&useIPv4, "t4", false, "Query and show IPv4 addresses")
    whoisCmd.Flags().BoolVar(&useIPv6, "t6", false, "Query and show IPv6 addresses")
    whoisCmd.Flags().BoolVar(&outputJSON, "json", false, "Output in JSON format")
    whoisCmd.Flags().BoolVar(&outputCSV, "csv", false, "Output in CSV format")
    whoisCmd.Flags().IntVarP(&workers, "workers", "w", 3, "Number of concurrent workers")
    whoisCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output mode")
    whoisCmd.Flags().IntVar(&retry, "retry", 2, "Number of retries on failure")
    whoisCmd.Flags().IntVar(&delay, "delay", 1000, "Delay between queries in milliseconds")

    rootCmd.AddCommand(whoisCmd)
}

func runWhois(cmd *cobra.Command, args []string) {
    var domains []string

    if domain != "" {
        domains = append(domains, domain)
    }

    if domainFile != "" {
        file, err := os.Open(domainFile)
        if err != nil {
            fmt.Printf("❌ Unable to read file: %v\n", err)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if line != "" {
                domains = append(domains, line)
            }
        }
    }

    if len(domains) == 0 {
        fmt.Println("❌ Please use -d or --domainfile to provide domains")
        return
    }

    // Use goroutines and channels for concurrent queries
    results := processDomainsWithWorkers(domains, workers)

    // Display results based on output format
    switch {
    case outputJSON:
        outputResultsAsJSON(results)
    case outputCSV:
        outputResultsAsCSV(results)
    default:
        outputResultsAsTable(results)
    }
}

// processDomainsWithWorkers uses worker pool pattern to process domain queries
func processDomainsWithWorkers(domains []string, numWorkers int) []WhoisResult {
    domainChan := make(chan string, len(domains))
    resultChan := make(chan WhoisResult, len(domains))
    var wg sync.WaitGroup

    // Start worker goroutines
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(domainChan, resultChan, &wg)
    }

    // Send all domains to channel
    for _, d := range domains {
        domainChan <- d
    }
    close(domainChan)

    // Wait for all work to complete
    go func() {
        wg.Wait()
        close(resultChan)
    }()

    // Collect results
    var results []WhoisResult
    for result := range resultChan {
        results = append(results, result)
    }

    return results
}

// worker is the worker function for querying individual domains
func worker(domainChan <-chan string, resultChan chan<- WhoisResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for d := range domainChan {
        if verbose {
            fmt.Fprintf(os.Stderr, "🔍 Querying %s...\n", d)
        }
        
        result := queryDomainWithRetry(d)
        resultChan <- result
        
        // Add query interval to avoid too frequent requests
        if delay > 0 {
            time.Sleep(time.Duration(delay) * time.Millisecond)
        }
    }
}

// queryDomainWithRetry queries WHOIS information for a single domain (with retry mechanism)
func queryDomainWithRetry(d string) WhoisResult {
    var lastErr error
    
    for attempt := 0; attempt <= retry; attempt++ {
        if attempt > 0 {
            if verbose {
                fmt.Fprintf(os.Stderr, "⚠️  %s retry attempt %d...\n", d, attempt)
            }
            // Wait longer before retry
            time.Sleep(time.Duration(delay*2) * time.Millisecond)
        }
        
        result := queryDomain(d)
        if result.Status == "success" {
            if verbose {
                fmt.Fprintf(os.Stderr, "✅ %s query successful\n", d)
            }
            return result
        }
        
        lastErr = fmt.Errorf(result.Error)
    }
    
    if verbose {
        fmt.Fprintf(os.Stderr, "❌ %s query failed (retried %d times): %v\n", d, retry, lastErr)
    }
    
    return WhoisResult{
        Domain:    d,
        Registrar: "Query Failed",
        Expires:   "-",
        DaysLeft:  "-",
        Status:    "failed",
        Error:     fmt.Sprintf("Failed after %d retries: %v", retry, lastErr),
    }
}

// queryDomain queries WHOIS information for a single domain
func queryDomain(d string) WhoisResult {
    result := WhoisResult{
        Domain:    d,
        Registrar: "-",
        Expires:   "-",
        DaysLeft:  "-",
        Status:    "success",
        Error:     "",
        ASN:       "-",
        IPv4:      "-",
        IPv6:      "-",
        Country:   "-",
    }

    raw, err := whois.Whois(d)
    if err != nil {
        result.Status = "failed"
        result.Registrar = "Query Failed"
        result.Error = fmt.Sprintf("WHOIS query failed: %v", err)
        return result
    }

    if strings.TrimSpace(raw) == "" {
        result.Status = "failed"
        result.Registrar = "Query Failed"
        result.Error = "WHOIS query returned empty result"
        return result
    }

    parsed, err := whoisparser.Parse(raw)
    if err != nil {
        result.Status = "failed"
        result.Registrar = "Parse Failed"
        result.Error = fmt.Sprintf("WHOIS parsing failed: %v", err)
        return result
    }

    // Basic domain information
    exp := parsed.Domain.ExpirationDate
    registrar := parsed.Registrar.Name
    
    // If ASN query is enabled, try to get ASN information
    if showASN {
        result.ASN = extractASNFromDomain(d)
    }
    
    // If IPv4/IPv6 query is enabled, try to resolve IP addresses
    if useIPv4 || useIPv6 {
        ipv4, ipv6 := extractIPAddresses(d)
        if useIPv4 {
            result.IPv4 = ipv4
        }
        if useIPv6 {
            result.IPv6 = ipv6
        }
    }
    
    // Try to get country information
    result.Country = extractCountryFromRaw(raw)
    
    if exp != "" {
        // Try multiple date formats
        layouts := []string{
            "2006-01-02T15:04:05Z",
            "2006-01-02T15:04:05.000Z",
            "2006-01-02 15:04:05",
            "2006-01-02",
            "2006-01-02 15:04:05 (UTC+8)",
            "2006-01-02 15:04:05 (MST)",
        }
        
        cleanExp := strings.TrimSpace(exp)
        var parsedTime time.Time
        var err error
        
        for _, layout := range layouts {
            if parsedTime, err = time.Parse(layout, cleanExp); err == nil {
                break
            }
        }
        
        if err != nil && strings.Contains(cleanExp, " ") {
            datePart := strings.Fields(cleanExp)[0]
            parsedTime, err = time.Parse("2006-01-02", datePart)
        }
        
        if err == nil {
            days := int(time.Until(parsedTime).Hours() / 24)
            result.DaysLeft = strconv.Itoa(days)
            result.Expires = parsedTime.Format("2006-01-02")
        } else {
            if len(cleanExp) >= 10 {
                result.Expires = cleanExp[:10]
            } else {
                result.Expires = cleanExp
            }
        }
    }

    if len(registrar) > 20 {
        registrar = registrar[:17] + "..."
    }
    result.Registrar = registrar

    return result
}

// extractASNFromDomain extracts ASN information by querying the domain's IP address
func extractASNFromDomain(domain string) string {
    // First, get the IP address of the domain
    ips, err := net.LookupIP(domain)
    if err != nil {
        return "-"
    }
    
    // Try to get ASN for the first IPv4 address
    for _, ip := range ips {
        if ip.To4() != nil {
            // Query ASN information for this IP
            asn := queryASNForIP(ip.String())
            if asn != "-" {
                return asn
            }
        }
    }
    
    return "-"
}

// queryASNForIP queries ASN information for a given IP address
func queryASNForIP(ip string) string {
    // Use a timeout context to prevent hanging
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // First try a quick DNS-based ASN lookup
    if asn := queryASNViaOriginAS(ip); asn != "-" {
        return asn
    }
    
    // If that fails, try WHOIS with timeout
    done := make(chan string, 1)
    go func() {
        raw, err := whois.Whois(ip)
        if err != nil {
            done <- "-"
            return
        }
        done <- extractASNFromRaw(raw)
    }()
    
    select {
    case result := <-done:
        return result
    case <-ctx.Done():
        return "-"
    }
}

// queryASNViaOriginAS queries ASN using origin.asn.cymru.com
func queryASNViaOriginAS(ip string) string {
    // Split IP into parts and reverse for origin.asn.cymru.com query
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return "-"
    }
    
    // Reverse the IP for the DNS query
    reversed := fmt.Sprintf("%s.%s.%s.%s.origin.asn.cymru.com", parts[3], parts[2], parts[1], parts[0])
    
    // Query TXT record
    txtRecords, err := net.LookupTXT(reversed)
    if err != nil {
        return "-"
    }
    
    // Parse the TXT record to extract ASN
    for _, record := range txtRecords {
        fields := strings.Fields(record)
        if len(fields) >= 1 {
            asn := strings.Trim(fields[0], "\"")
            if _, err := strconv.Atoi(asn); err == nil {
                return "AS" + asn
            }
        }
    }
    
    return "-"
}

// extractASNFromRaw extracts ASN information from raw WHOIS data
func extractASNFromRaw(raw string) string {
    lines := strings.Split(raw, "\n")
    for _, line := range lines {
        originalLine := strings.TrimSpace(line)
        line = strings.ToLower(originalLine)
        
        // Look for various ASN patterns
        patterns := []string{
            "asn:",
            "as number:",
            "autonomous system:",
            "origin as:",
            "originas:",
            "origin:",
            "asnum:",
            "aut-num:",
        }
        
        for _, pattern := range patterns {
            if strings.Contains(line, pattern) {
                // Extract the part after the pattern
                parts := strings.Split(originalLine, ":")
                if len(parts) >= 2 {
                    value := strings.TrimSpace(parts[1])
                    // Look for AS followed by numbers
                    fields := strings.Fields(value)
                    for _, field := range fields {
                        if matched := extractASNumber(field); matched != "" {
                            return matched
                        }
                    }
                }
                
                // Also check all fields in the line
                fields := strings.Fields(originalLine)
                for _, field := range fields {
                    if matched := extractASNumber(field); matched != "" {
                        return matched
                    }
                }
            }
        }
    }
    return "-"
}

// extractASNumber extracts AS number from a string field
func extractASNumber(field string) string {
    field = strings.TrimSpace(field)
    
    // Remove common separators and prefixes
    field = strings.Trim(field, "(),[]\"'")
    
    // Check if it starts with AS followed by numbers
    if strings.HasPrefix(strings.ToUpper(field), "AS") && len(field) > 2 {
        asNum := field[2:]
        // Check if the rest are digits
        if _, err := strconv.Atoi(asNum); err == nil {
            return "AS" + asNum
        }
    }
    
    // Check if it's just a number that could be an ASN
    if num, err := strconv.Atoi(field); err == nil && num > 0 && num < 400000 {
        return "AS" + field
    }
    
    return ""
}

// extractIPAddresses gets IPv4 and IPv6 addresses for a domain
func extractIPAddresses(domain string) (string, string) {
    var ipv4, ipv6 string = "-", "-"
    
    // Query A records (IPv4)
    ips, err := net.LookupIP(domain)
    if err == nil {
        var ipv4List, ipv6List []string
        for _, ip := range ips {
            if ip.To4() != nil {
                ipv4List = append(ipv4List, ip.String())
            } else {
                ipv6List = append(ipv6List, ip.String())
            }
        }
        if len(ipv4List) > 0 {
            ipv4 = strings.Join(ipv4List, ",")
        }
        if len(ipv6List) > 0 {
            ipv6 = strings.Join(ipv6List, ",")
        }
    }
    
    return ipv4, ipv6
}

// extractCountryFromRaw extracts country information from raw WHOIS data
func extractCountryFromRaw(raw string) string {
    lines := strings.Split(raw, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(strings.ToLower(line))
        if strings.Contains(line, "country:") {
            parts := strings.Split(line, ":")
            if len(parts) >= 2 {
                country := strings.TrimSpace(parts[1])
                if len(country) > 0 {
                    return strings.ToUpper(country)
                }
            }
        }
        if strings.Contains(line, "registrant country:") {
            parts := strings.Split(line, ":")
            if len(parts) >= 2 {
                country := strings.TrimSpace(parts[1])
                if len(country) > 0 {
                    return strings.ToUpper(country)
                }
            }
        }
    }
    return "-"
}

// outputResultsAsJSON outputs results in JSON format
func outputResultsAsJSON(results []WhoisResult) {
    jsonData, err := json.MarshalIndent(results, "", "  ")
    if err != nil {
        fmt.Printf("❌ JSON encoding failed: %v\n", err)
        return
    }
    fmt.Println(string(jsonData))
}

// outputResultsAsCSV outputs results in CSV format
func outputResultsAsCSV(results []WhoisResult) {
    writer := csv.NewWriter(os.Stdout)
    defer writer.Flush()

    // Dynamically generate header row
    header := []string{"Domain", "Registrar", "Expires", "Days Left", "Status"}
    if showASN {
        header = append(header, "ASN")
    }
    if useIPv4 {
        header = append(header, "IPv4")
    }
    if useIPv6 {
        header = append(header, "IPv6")
    }
    header = append(header, "Country")
    
    if err := writer.Write(header); err != nil {
        fmt.Printf("❌ CSV write failed: %v\n", err)
        return
    }

    // Write data rows
    for _, result := range results {
        record := []string{
            result.Domain,
            result.Registrar,
            result.Expires,
            result.DaysLeft,
            result.Status,
        }
        if showASN {
            record = append(record, result.ASN)
        }
        if useIPv4 {
            record = append(record, result.IPv4)
        }
        if useIPv6 {
            record = append(record, result.IPv6)
        }
        record = append(record, result.Country)
        
        if err := writer.Write(record); err != nil {
            fmt.Printf("❌ CSV write failed: %v\n", err)
            return
        }
    }
}

// outputResultsAsTable outputs results in table format
func outputResultsAsTable(results []WhoisResult) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    
    // Dynamically adjust table headers based on enabled options
    headers := []interface{}{"Domain", "Registrar", "Expires", "Days Left"}
    if showASN {
        headers = append(headers, "ASN")
    }
    if useIPv4 {
        headers = append(headers, "IPv4")
    }
    if useIPv6 {
        headers = append(headers, "IPv6")
    }
    
    t.AppendHeader(table.Row(headers))

    for _, result := range results {
        row := []interface{}{result.Domain, result.Registrar, result.Expires, result.DaysLeft}
        if showASN {
            row = append(row, result.ASN)
        }
        if useIPv4 {
            row = append(row, result.IPv4)
        }
        if useIPv6 {
            row = append(row, result.IPv6)
        }
        t.AppendRow(table.Row(row))
    }

    t.SetStyle(table.StyleDefault)
    t.Render()
}