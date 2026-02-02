package database

import (
	"database/sql"
	"log"
	"net"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Attempt to resolve hostname to IPv4 to avoid IPv6 issues
	resolvedConnString := connectionString
	u, err := url.Parse(connectionString)
	if err == nil && u.Host != "" {
		host := u.Hostname()
		port := u.Port()

		// Only resolve if it's not already an IP
		if net.ParseIP(host) == nil {
			ips, err := net.LookupIP(host)
			if err == nil {
				var ipv4 net.IP
				for _, ip := range ips {
					if ip.To4() != nil {
						ipv4 = ip
						break
					}
				}
				if ipv4 != nil {
					log.Printf("Resolved database host %s to %s", host, ipv4.String())
					if port != "" {
						u.Host = net.JoinHostPort(ipv4.String(), port)
					} else {
						u.Host = ipv4.String()
					}
					resolvedConnString = u.String()
				}
			}
		}
	}

	// Open database
	db, err := sql.Open("pgx", resolvedConnString+"&prefer_simple_protocol=true")
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Println("Database connection failed (will retry):", err)
		return db, err // Return db object anyway so we can keep trying
	}

	// Set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
