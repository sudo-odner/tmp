package server

import "strings"

type Storage struct {
	data map[string]string
}

func New() *Storage {
	storage := map[string]string{
		"HTTP":  "80/TCP Hypertext Transfer Protocol unencrypted web traffic",
		"HTTPS": "443/TCP HTTP Secure encrypted web traffic TLS",
		"SSH":   "22/TCP Secure Shell remote administration",
		"FTP":   "21/TCP File Transfer Protocol control connection",
		"SMTP":  "25/TCP Simple Mail Transfer Protocol email transmission",
		"DHCP":  "67/UDP Dynamic Host Configuration Protocol IP assignment",
		"NTP":   "123/UDP Network Time Protocol clock synchronization",
		"BGP":   "179/TCP Border Gateway Protocol internet routing",
		"SNMP":  "161/UDP Simple Network Management Protocol device monitoring",
	}
	return &Storage{
		data: storage,
	}
}

func (s *Storage) ByKey(key string) (string, bool) {
	if val, ok := s.data[key]; ok {
		return val, true
	}
	val, ok := s.data[strings.ToUpper(key)]
	return val, ok
}

func (s *Storage) Search(substring string) [][2]string {
	list := make([][2]string, 0, len(s.data))
	subLower := strings.ToLower(substring)
	for key, value := range s.data {
		if strings.Contains(strings.ToLower(value), subLower) || strings.Contains(strings.ToLower(key), subLower) {
			list = append(list, [2]string{key, value})
		}
	}
	return list
}

func (s *Storage) List() []string {
	list := make([]string, 0, len(s.data))
	for key := range s.data {
		list = append(list, key)
	}
	return list
}

func (s *Storage) Count() int { return len(s.data) }
