package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"sort"
	"strings"
	"time"
)

type BitwardenExport struct {
	Encrypted bool              `json:"encrypted"`
	Items     []BitwardenItem   `json:"items"`
	Folders   []BitwardenFolder `json:"folders"`
}

type BitwardenFolder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BitwardenItem struct {
	ID               string              `json:"id"`
	OrganizationID   *string             `json:"organizationId"`
	FolderID         *string             `json:"folderId"`
	Type             int                 `json:"type"`
	Name             string              `json:"name"`
	Notes            *string             `json:"notes"`
	Favorite         bool                `json:"favorite"`
	Fields           []BitwardenField    `json:"fields"`
	Login            *BitwardenLogin     `json:"login"`
	Card             *BitwardenCard      `json:"card"`
	Identity         *BitwardenIdentity  `json:"identity"`
	SecureNote       *BitwardenSecureNote `json:"secureNote"`
	CreationDate     *string             `json:"creationDate"`
	RevisionDate     *string             `json:"revisionDate"`
}

type BitwardenLogin struct {
	Username            *string `json:"username"`
	Password            *string `json:"password"`
	PasswordRevisionDate *string `json:"passwordRevisionDate"`
	TOTP                *string `json:"totp"`
	URIs                []struct {
		Match *int    `json:"match"`
		URI   *string `json:"uri"`
	} `json:"uris"`
}

type BitwardenCard struct {
	CardholderName *string `json:"cardholderName"`
	Brand          *string `json:"brand"`
	Number         *string `json:"number"`
	ExpMonth       *string `json:"expMonth"`
	ExpYear        *string `json:"expYear"`
	Code           *string `json:"code"`
}

type BitwardenIdentity struct {
	Title      *string `json:"title"`
	FirstName  *string `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   *string `json:"lastName"`
	Address1   *string `json:"address1"`
	Address2   *string `json:"address2"`
	Address3   *string `json:"address3"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postalCode"`
	Country    *string `json:"country"`
	Company    *string `json:"company"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	SSN        *string `json:"ssn"`
	Username   *string `json:"username"`
	PassportNumber *string `json:"passportNumber"`
	LicenseNumber  *string `json:"licenseNumber"`
}

type BitwardenSecureNote struct {
	Type int `json:"type"`
}

type BitwardenField struct {
	Name  string  `json:"name"`
	Value *string `json:"value"`
	Type  int     `json:"type"`
}

type GroupingMode int

const (
	GroupNone GroupingMode = iota
	GroupByType
	GroupByFolder
)

type ExportFields struct {
	Type              bool
	Name              bool
	Username          bool
	Password          bool
	Notes             bool
	URL               bool
	Favorite          bool
	TOTP              bool
	CreationDate      bool
	ModificationDate  bool
	PasswordRevision  bool
	Folder            bool
	Organization      bool
	CustomFields      bool
}

type ExportOptions struct {
	Fields   ExportFields
	Grouping GroupingMode
}

func ConvertBitwardenToHTML(inputPath, outputPath string, options ExportOptions) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var export BitwardenExport
	if err := json.Unmarshal(data, &export); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	if export.Encrypted {
		return fmt.Errorf("file is encrypted, please export as unencrypted JSON")
	}

	htmlContent := generateHTML(export, options)

	if err := os.WriteFile(outputPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("error writing HTML file: %w", err)
	}

	return nil
}

func generateHTML(export BitwardenExport, options ExportOptions) string {
	var sb strings.Builder
	fields := options.Fields

	// Create folder map for resolving folder names
	folderMap := make(map[string]string)
	for _, folder := range export.Folders {
		folderMap[folder.ID] = folder.Name
	}

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bitwarden Export</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            padding: 20px;
            background: #f5f5f5;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            position: relative;
        }

        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }

        .export-info {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }

        .section-header {
            background: #f8f9fa;
            padding: 12px 16px;
            margin-top: 20px;
            margin-bottom: 10px;
            border-radius: 6px;
            cursor: pointer;
            user-select: none;
            display: flex;
            align-items: center;
            justify-content: space-between;
            transition: background 0.2s;
        }

        .section-header:hover {
            background: #e9ecef;
        }

        .section-header h2 {
            font-size: 18px;
            color: #175ddc;
            margin: 0;
        }

        .section-toggle {
            font-size: 20px;
            font-weight: bold;
            color: #175ddc;
            transition: transform 0.2s;
        }

        .section-header.collapsed .section-toggle {
            transform: rotate(-90deg);
        }

        .section-content {
            display: block;
            overflow: hidden;
            transition: max-height 0.3s ease-out;
        }

        .section-content.collapsed {
            display: none;
        }

        .section-count {
            font-size: 14px;
            color: #666;
            margin-left: 10px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }

        th {
            background: #175ddc;
            color: white;
            padding: 12px 8px;
            text-align: left;
            font-weight: 600;
            font-size: 13px;
            position: sticky;
            top: 0;
        }

        th.sortable {
            cursor: pointer;
            user-select: none;
        }

        th.sortable:hover {
            background: #1450bc;
        }

        th.sortable::after {
            content: ' ⇅';
            opacity: 0.5;
        }

        th.sortable.asc::after {
            content: ' ↑';
            opacity: 1;
        }

        th.sortable.desc::after {
            content: ' ↓';
            opacity: 1;
        }

        td {
            padding: 10px 8px;
            border-bottom: 1px solid #e0e0e0;
            font-size: 13px;
            vertical-align: top;
        }

        tr:nth-child(even) {
            background: #f8f9fa;
        }

        tr:hover {
            background: #e3f2fd;
        }

        .item-type {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
        }

        .type-login { background: #e3f2fd; color: #1976d2; }
        .type-card { background: #fff3e0; color: #f57c00; }
        .type-identity { background: #f3e5f5; color: #7b1fa2; }
        .type-note { background: #e8f5e9; color: #388e3c; }

        .password {
            font-family: 'Courier New', monospace;
            background: #f5f5f5;
            padding: 2px 6px;
            border-radius: 3px;
            word-break: break-all;
        }

        .url {
            color: #175ddc;
            text-decoration: none;
            word-break: break-all;
        }

        .notes {
            max-width: 300px;
            white-space: pre-wrap;
            word-break: break-word;
        }

        .favorite {
            color: #ffd700;
            font-size: 16px;
        }

        .empty {
            color: #999;
            font-style: italic;
        }

        .date {
            font-size: 11px;
            color: #666;
        }

        .search-container {
            margin: 20px 0;
        }

        .search-box {
            width: 100%;
            padding: 12px 16px;
            font-size: 14px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            outline: none;
            transition: border-color 0.2s;
        }

        .search-box:focus {
            border-color: #175ddc;
        }

        .search-info {
            margin-top: 8px;
            font-size: 13px;
            color: #666;
        }

        .print-icon {
            position: absolute;
            top: 30px;
            right: 30px;
            width: 32px;
            height: 32px;
            cursor: pointer;
            opacity: 0.7;
            transition: opacity 0.2s;
        }

        .print-icon:hover {
            opacity: 1;
        }

        .print-icon svg {
            width: 100%;
            height: 100%;
        }

        .toggle-password-btn {
            position: absolute;
            top: 30px;
            right: 80px;
            width: 32px;
            height: 32px;
            cursor: pointer;
            opacity: 0.7;
            transition: opacity 0.2s;
            background: none;
            border: none;
            padding: 0;
        }

        .toggle-password-btn:hover {
            opacity: 1;
        }

        .toggle-password-btn svg {
            width: 100%;
            height: 100%;
        }

        .password.masked {
            user-select: none;
        }

        .controls-container {
            margin: 20px 0;
        }

        .search-wrapper {
            width: 100%;
        }

        tr.hidden {
            display: none;
        }

        @media print {
            body {
                background: white;
                padding: 0;
            }

            .container {
                box-shadow: none;
                padding: 0;
                max-width: 100%;
            }

            .search-container,
            .controls-container,
            .print-button,
            .print-icon {
                display: none;
            }

            .export-info {
                margin-bottom: 10px;
                font-size: 10px;
            }

            h1 {
                font-size: 18px;
                margin-bottom: 5px;
            }

            table {
                page-break-inside: auto;
                font-size: 9px;
            }

            thead {
                display: table-header-group;
            }

            tr {
                page-break-inside: avoid;
                page-break-after: auto;
            }

            th {
                position: static;
                padding: 6px 4px;
                font-size: 9px;
            }

            td {
                padding: 4px 4px;
                font-size: 9px;
            }

            tr:hover {
                background: transparent;
            }

            .password,
            .notes {
                font-size: 8px;
                max-width: 150px;
            }

            .item-type {
                font-size: 8px;
                padding: 2px 4px;
            }

            .date {
                font-size: 8px;
            }

            .url {
                font-size: 8px;
                word-break: break-all;
            }

            /* Hide rows that are filtered out */
            tr.hidden {
                display: none;
            }

            @page {
                margin: 0.5cm;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Bitwarden Export</h1>
        <div class="export-info">
            Generated on ` + time.Now().Format("01/02/2006 at 3:04:05 PM") + `<br>
            Number of entries: ` + fmt.Sprintf("%d", len(export.Items)) + `
        </div>

        <div class="print-icon" onclick="window.print()" title="Print">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#175ddc" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="6 9 6 2 18 2 18 9"></polyline>
                <path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"></path>
                <rect x="6" y="14" width="12" height="8"></rect>
            </svg>
        </div>

        <button class="toggle-password-btn" id="togglePasswordBtn" title="Toggle password visibility">
            <svg id="eyeIcon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#175ddc" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                <circle cx="12" cy="12" r="3"></circle>
            </svg>
        </button>

        <div class="controls-container">
            <div class="search-wrapper">
                <input type="text" class="search-box" id="searchInput" placeholder="Search by name, username, URL, notes...">
                <div class="search-info">
                    <span id="searchResults"></span>
                </div>
            </div>
        </div>

        <table>
            <thead>
                <tr>
`)

	// Generate table headers based on selected fields
	if fields.Type {
		sb.WriteString("                    <th>Type</th>\n")
	}
	if fields.Name {
		sb.WriteString("                    <th>Name</th>\n")
	}
	if fields.Username {
		sb.WriteString("                    <th>Username</th>\n")
	}
	if fields.Password {
		sb.WriteString("                    <th>Password</th>\n")
	}
	if fields.URL {
		sb.WriteString("                    <th>URL</th>\n")
	}
	if fields.Notes {
		sb.WriteString("                    <th>Notes</th>\n")
	}
	if fields.TOTP {
		sb.WriteString("                    <th>TOTP</th>\n")
	}
	if fields.Folder {
		sb.WriteString("                    <th>Folder</th>\n")
	}
	if fields.Organization {
		sb.WriteString("                    <th>Organization</th>\n")
	}
	if fields.CreationDate {
		sb.WriteString("                    <th>Created</th>\n")
	}
	if fields.ModificationDate {
		sb.WriteString("                    <th>Modified</th>\n")
	}
	if fields.PasswordRevision {
		sb.WriteString("                    <th>Password Revision</th>\n")
	}
	if fields.CustomFields {
		sb.WriteString("                    <th>Custom Fields</th>\n")
	}

	sb.WriteString(`                </tr>
            </thead>
            <tbody>
`)

	// Generate rows based on grouping mode
	switch options.Grouping {
	case GroupByType:
		sb.WriteString(`            </tbody>
        </table>
`)
		groups := groupItemsByType(export.Items)
		groupNames := []string{"Login", "Note", "Card", "Identity", "Other"}
		for _, groupName := range groupNames {
			items, exists := groups[groupName]
			if !exists || len(items) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf(`        <div class="section-header" onclick="toggleSection(this)">
            <div>
                <h2>%s<span class="section-count">(%d entries)</span></h2>
            </div>
            <span class="section-toggle">▼</span>
        </div>
        <div class="section-content">
        <table>
            <thead>
                <tr>
`, groupName, len(items)))
			// Repeat headers for grouped tables
			generateTableHeaders(&sb, fields)
			sb.WriteString(`                </tr>
            </thead>
            <tbody>
`)
			for _, item := range items {
				generateItemRow(&sb, item, fields, folderMap)
			}
			sb.WriteString(`            </tbody>
        </table>
        </div>
`)
		}

	case GroupByFolder:
		sb.WriteString(`            </tbody>
        </table>
`)
		groups := groupItemsByFolder(export.Items, export.Folders)
		var folderNames []string
		for folderName := range groups {
			folderNames = append(folderNames, folderName)
		}
		// Sort: "No Folder" first, then alphabetically
		var sortedNames []string
		hasNoFolder := false
		for _, name := range folderNames {
			if name == "No Folder" {
				hasNoFolder = true
			} else {
				sortedNames = append(sortedNames, name)
			}
		}
		sort.Strings(sortedNames)
		if hasNoFolder {
			folderNames = append([]string{"No Folder"}, sortedNames...)
		} else {
			folderNames = sortedNames
		}

		for _, folderName := range folderNames {
			items := groups[folderName]
			sb.WriteString(fmt.Sprintf(`        <div class="section-header" onclick="toggleSection(this)">
            <div>
                <h2>📁 %s<span class="section-count">(%d entries)</span></h2>
            </div>
            <span class="section-toggle">▼</span>
        </div>
        <div class="section-content">
        <table>
            <thead>
                <tr>
`, html.EscapeString(folderName), len(items)))
			generateTableHeaders(&sb, fields)
			sb.WriteString(`                </tr>
            </thead>
            <tbody>
`)
			for _, item := range items {
				generateItemRow(&sb, item, fields, folderMap)
			}
			sb.WriteString(`            </tbody>
        </table>
        </div>
`)
		}

	default:
		// No grouping - original behavior
		for _, item := range export.Items {
			generateItemRow(&sb, item, fields, folderMap)
		}
		sb.WriteString(`            </tbody>
        </table>
`)
	}

	sb.WriteString(`    </div>
    <script>
        function toggleSection(header) {
            const content = header.nextElementSibling;

            if (!content || !content.classList.contains('section-content')) {
                console.error('Section content not found');
                return;
            }

            // Simple toggle with display
            header.classList.toggle('collapsed');
            content.classList.toggle('collapsed');
        }

        document.addEventListener('DOMContentLoaded', function() {

            const tables = document.querySelectorAll('table');
            const searchInput = document.getElementById('searchInput');
            const searchResults = document.getElementById('searchResults');

            let allRows = [];
            tables.forEach(table => {
                const tbody = table.querySelector('tbody');
                if (tbody) {
                    allRows = allRows.concat(Array.from(tbody.querySelectorAll('tr')));
                }
            });
            const totalEntries = allRows.length;

            // Password toggle functionality
            const togglePasswordBtn = document.getElementById('togglePasswordBtn');
            const eyeIcon = document.getElementById('eyeIcon');
            let passwordsVisible = true;

            if (togglePasswordBtn) {
                togglePasswordBtn.addEventListener('click', function() {
                    passwordsVisible = !passwordsVisible;
                    const passwordCells = document.querySelectorAll('.password');

                    passwordCells.forEach(cell => {
                        if (passwordsVisible) {
                            cell.classList.remove('masked');
                            if (cell.dataset.original) {
                                cell.textContent = cell.dataset.original;
                            }
                        } else {
                            if (!cell.dataset.original) {
                                cell.dataset.original = cell.textContent;
                            }
                            const length = cell.textContent.length;
                            cell.textContent = '•'.repeat(Math.min(length, 20));
                            cell.classList.add('masked');
                        }
                    });

                    // Update icon
                    if (passwordsVisible) {
                        eyeIcon.innerHTML = '<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle>';
                    } else {
                        eyeIcon.innerHTML = '<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line>';
                    }
                });
            }

            // Search/Filter functionality
            searchInput.addEventListener('input', function() {
                const searchTerm = this.value.toLowerCase().trim();
                let visibleCount = 0;

                allRows.forEach(row => {
                    if (searchTerm === '') {
                        row.classList.remove('hidden');
                        visibleCount++;
                    } else {
                        const searchableText = row.textContent.toLowerCase();
                        if (searchableText.includes(searchTerm)) {
                            row.classList.remove('hidden');
                            visibleCount++;
                        } else {
                            row.classList.add('hidden');
                        }
                    }
                });

                // Update search results info
                if (searchTerm === '') {
                    searchResults.textContent = '';
                } else {
                    searchResults.textContent = 'Showing ' + visibleCount + ' of ' + totalEntries + ' entries';
                }
            });

            // Sorting functionality for each table
            tables.forEach(table => {
                const headers = table.querySelectorAll('th');
                const tbody = table.querySelector('tbody');
                if (!tbody) return;

                headers.forEach((header, index) => {
                    // Don't make password column sortable
                    const headerText = header.textContent.trim();
                    if (headerText.toLowerCase().includes('password') && !headerText.toLowerCase().includes('revision')) {
                        return;
                    }

                    header.classList.add('sortable');
                    let ascending = true;

                    header.addEventListener('click', function() {
                        // Remove sort indicators from other headers
                        headers.forEach(h => {
                            h.classList.remove('asc', 'desc');
                        });

                        const rows = Array.from(tbody.querySelectorAll('tr'));

                    rows.sort((a, b) => {
                        const aCell = a.children[index];
                        const bCell = b.children[index];

                        if (!aCell || !bCell) return 0;

                        let aValue = aCell.textContent.trim();
                        let bValue = bCell.textContent.trim();

                        // Handle empty values
                        if (aValue === '-') aValue = '';
                        if (bValue === '-') bValue = '';

                        // Try to parse as date
                        const aDate = new Date(aValue);
                        const bDate = new Date(bValue);

                        if (!isNaN(aDate) && !isNaN(bDate) && aValue.match(/\d{4}-\d{2}-\d{2}/)) {
                            return ascending ? aDate - bDate : bDate - aDate;
                        }

                        // Compare as strings (case insensitive)
                        const comparison = aValue.toLowerCase().localeCompare(bValue.toLowerCase());
                        return ascending ? comparison : -comparison;
                    });

                    // Update table
                    rows.forEach(row => tbody.appendChild(row));

                    // Update sort indicator
                    header.classList.add(ascending ? 'asc' : 'desc');
                    ascending = !ascending;
                    });
                });
            });
        });
    </script>
</body>
</html>`)

	return sb.String()
}

func getItemType(typeID int) (class, name string) {
	switch typeID {
	case 1:
		return "type-login", "Login"
	case 2:
		return "type-note", "Note"
	case 3:
		return "type-card", "Card"
	case 4:
		return "type-identity", "Identity"
	default:
		return "type-note", "Other"
	}
}

func getUsername(item BitwardenItem) string {
	if item.Login != nil && item.Login.Username != nil {
		return *item.Login.Username
	}
	if item.Identity != nil && item.Identity.Username != nil {
		return *item.Identity.Username
	}
	if item.Identity != nil && item.Identity.Email != nil {
		return *item.Identity.Email
	}
	return ""
}

func getPassword(item BitwardenItem) string {
	if item.Login != nil && item.Login.Password != nil {
		return *item.Login.Password
	}
	if item.Card != nil && item.Card.Code != nil {
		return "CVV: " + *item.Card.Code
	}
	return ""
}

func getURLs(item BitwardenItem) []string {
	var urls []string
	if item.Login != nil {
		for _, uri := range item.Login.URIs {
			if uri.URI != nil && *uri.URI != "" {
				urls = append(urls, *uri.URI)
			}
		}
	}
	return urls
}

func getNotes(item BitwardenItem) string {
	if item.Notes != nil {
		return *item.Notes
	}

	// For cards, add info
	if item.Card != nil {
		var parts []string
		if item.Card.CardholderName != nil {
			parts = append(parts, "Cardholder: "+*item.Card.CardholderName)
		}
		if item.Card.Brand != nil {
			parts = append(parts, "Brand: "+*item.Card.Brand)
		}
		if item.Card.Number != nil {
			parts = append(parts, "Number: "+*item.Card.Number)
		}
		if item.Card.ExpMonth != nil && item.Card.ExpYear != nil {
			parts = append(parts, fmt.Sprintf("Expiration: %s/%s", *item.Card.ExpMonth, *item.Card.ExpYear))
		}
		return strings.Join(parts, "\n")
	}

	// For identities
	if item.Identity != nil {
		var parts []string
		if item.Identity.FirstName != nil || item.Identity.LastName != nil {
			name := ""
			if item.Identity.FirstName != nil {
				name += *item.Identity.FirstName
			}
			if item.Identity.LastName != nil {
				if name != "" {
					name += " "
				}
				name += *item.Identity.LastName
			}
			parts = append(parts, "Name: "+name)
		}
		if item.Identity.Company != nil {
			parts = append(parts, "Company: "+*item.Identity.Company)
		}
		if item.Identity.Phone != nil {
			parts = append(parts, "Phone: "+*item.Identity.Phone)
		}
		if item.Identity.Address1 != nil {
			parts = append(parts, "Address: "+*item.Identity.Address1)
		}
		if item.Identity.City != nil || item.Identity.PostalCode != nil {
			city := ""
			if item.Identity.PostalCode != nil {
				city = *item.Identity.PostalCode
			}
			if item.Identity.City != nil {
				if city != "" {
					city += " "
				}
				city += *item.Identity.City
			}
			if city != "" {
				parts = append(parts, city)
			}
		}
		return strings.Join(parts, "\n")
	}

	return ""
}

func getTOTP(item BitwardenItem) string {
	if item.Login != nil && item.Login.TOTP != nil {
		return *item.Login.TOTP
	}
	return ""
}

func getPasswordRevisionDate(item BitwardenItem) string {
	if item.Login != nil && item.Login.PasswordRevisionDate != nil {
		return *item.Login.PasswordRevisionDate
	}
	return ""
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02 15:04")
}

func getCustomFields(item BitwardenItem) string {
	if len(item.Fields) == 0 {
		return ""
	}

	var parts []string
	for _, field := range item.Fields {
		if field.Value != nil && *field.Value != "" {
			fieldType := ""
			switch field.Type {
			case 1:
				fieldType = " (hidden)"
			case 2:
				fieldType = " (boolean)"
			}
			parts = append(parts, fmt.Sprintf("%s: %s%s", field.Name, *field.Value, fieldType))
		}
	}
	return strings.Join(parts, "\n")
}

func generateTableHeaders(sb *strings.Builder, fields ExportFields) {
	if fields.Type {
		sb.WriteString("                    <th>Type</th>\n")
	}
	if fields.Name {
		sb.WriteString("                    <th>Name</th>\n")
	}
	if fields.Username {
		sb.WriteString("                    <th>Username</th>\n")
	}
	if fields.Password {
		sb.WriteString("                    <th>Password</th>\n")
	}
	if fields.URL {
		sb.WriteString("                    <th>URL</th>\n")
	}
	if fields.Notes {
		sb.WriteString("                    <th>Notes</th>\n")
	}
	if fields.TOTP {
		sb.WriteString("                    <th>TOTP</th>\n")
	}
	if fields.Folder {
		sb.WriteString("                    <th>Folder</th>\n")
	}
	if fields.Organization {
		sb.WriteString("                    <th>Organization</th>\n")
	}
	if fields.CreationDate {
		sb.WriteString("                    <th>Created</th>\n")
	}
	if fields.ModificationDate {
		sb.WriteString("                    <th>Modified</th>\n")
	}
	if fields.PasswordRevision {
		sb.WriteString("                    <th>Password Revision</th>\n")
	}
	if fields.CustomFields {
		sb.WriteString("                    <th>Custom Fields</th>\n")
	}
}

func groupItemsByType(items []BitwardenItem) map[string][]BitwardenItem {
	groups := make(map[string][]BitwardenItem)
	for _, item := range items {
		_, typeName := getItemType(item.Type)
		groups[typeName] = append(groups[typeName], item)
	}
	return groups
}

func groupItemsByFolder(items []BitwardenItem, folders []BitwardenFolder) map[string][]BitwardenItem {
	// Create a map for quick folder name lookup
	folderMap := make(map[string]string)
	for _, folder := range folders {
		folderMap[folder.ID] = folder.Name
	}

	groups := make(map[string][]BitwardenItem)
	for _, item := range items {
		folderName := "No Folder"
		if item.FolderID != nil && *item.FolderID != "" {
			if name, exists := folderMap[*item.FolderID]; exists {
				folderName = name
			} else {
				folderName = *item.FolderID // Fallback to ID if name not found
			}
		}
		groups[folderName] = append(groups[folderName], item)
	}
	return groups
}

func generateItemRow(sb *strings.Builder, item BitwardenItem, fields ExportFields, folderMap map[string]string) {
	sb.WriteString("                <tr>\n")

	// Type
	if fields.Type {
		sb.WriteString("                    <td>")
		typeClass, typeName := getItemType(item.Type)
		if fields.Favorite && item.Favorite {
			sb.WriteString("<span class=\"favorite\">★</span> ")
		}
		sb.WriteString(fmt.Sprintf("<span class=\"item-type %s\">%s</span>", typeClass, typeName))
		sb.WriteString("</td>\n")
	}

	// Name
	if fields.Name {
		sb.WriteString("                    <td><strong>")
		sb.WriteString(html.EscapeString(item.Name))
		sb.WriteString("</strong></td>\n")
	}

	// Username/Login
	if fields.Username {
		sb.WriteString("                    <td>")
		username := getUsername(item)
		if username != "" {
			sb.WriteString(html.EscapeString(username))
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Password
	if fields.Password {
		sb.WriteString("                    <td>")
		password := getPassword(item)
		if password != "" {
			sb.WriteString("<span class=\"password\">")
			sb.WriteString(html.EscapeString(password))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// URL
	if fields.URL {
		sb.WriteString("                    <td>")
		urls := getURLs(item)
		if len(urls) > 0 {
			for i, url := range urls {
				if i > 0 {
					sb.WriteString("<br>")
				}
				sb.WriteString(fmt.Sprintf("<a href=\"%s\" class=\"url\">%s</a>",
					html.EscapeString(url), html.EscapeString(url)))
			}
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Notes
	if fields.Notes {
		sb.WriteString("                    <td>")
		notes := getNotes(item)
		if notes != "" {
			sb.WriteString("<span class=\"notes\">")
			sb.WriteString(html.EscapeString(notes))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// TOTP
	if fields.TOTP {
		sb.WriteString("                    <td>")
		totp := getTOTP(item)
		if totp != "" {
			sb.WriteString("<span class=\"password\">")
			sb.WriteString(html.EscapeString(totp))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Folder
	if fields.Folder {
		sb.WriteString("                    <td>")
		if item.FolderID != nil && *item.FolderID != "" {
			folderName := *item.FolderID
			if name, exists := folderMap[*item.FolderID]; exists {
				folderName = name
			}
			sb.WriteString(html.EscapeString(folderName))
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Organization
	if fields.Organization {
		sb.WriteString("                    <td>")
		if item.OrganizationID != nil && *item.OrganizationID != "" {
			sb.WriteString(html.EscapeString(*item.OrganizationID))
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Creation Date
	if fields.CreationDate {
		sb.WriteString("                    <td>")
		if item.CreationDate != nil && *item.CreationDate != "" {
			sb.WriteString("<span class=\"date\">")
			sb.WriteString(formatDate(*item.CreationDate))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Modification Date
	if fields.ModificationDate {
		sb.WriteString("                    <td>")
		if item.RevisionDate != nil && *item.RevisionDate != "" {
			sb.WriteString("<span class=\"date\">")
			sb.WriteString(formatDate(*item.RevisionDate))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Password Revision
	if fields.PasswordRevision {
		sb.WriteString("                    <td>")
		passwordRevision := getPasswordRevisionDate(item)
		if passwordRevision != "" {
			sb.WriteString("<span class=\"date\">")
			sb.WriteString(formatDate(passwordRevision))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	// Custom Fields
	if fields.CustomFields {
		sb.WriteString("                    <td>")
		customFields := getCustomFields(item)
		if customFields != "" {
			sb.WriteString("<span class=\"notes\">")
			sb.WriteString(html.EscapeString(customFields))
			sb.WriteString("</span>")
		} else {
			sb.WriteString("<span class=\"empty\">-</span>")
		}
		sb.WriteString("</td>\n")
	}

	sb.WriteString("                </tr>\n")
}
