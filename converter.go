package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
	"time"
)

type BitwardenExport struct {
	Encrypted bool            `json:"encrypted"`
	Items     []BitwardenItem `json:"items"`
}

type BitwardenItem struct {
	ID               string              `json:"id"`
	OrganizationID   *string             `json:"organizationId"`
	FolderID         *string             `json:"folderId"`
	Type             int                 `json:"type"`
	Name             string              `json:"name"`
	Notes            *string             `json:"notes"`
	Favorite         bool                `json:"favorite"`
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
}

func ConvertBitwardenToHTML(inputPath, outputPath string, fields ExportFields) error {
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

	htmlContent := generateHTML(export, fields)

	if err := os.WriteFile(outputPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("error writing HTML file: %w", err)
	}

	return nil
}

func generateHTML(export BitwardenExport, fields ExportFields) string {
	var sb strings.Builder

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

        tr:hover {
            background: #f8f9fa;
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

        @media print {
            body {
                background: white;
                padding: 0;
            }

            .container {
                box-shadow: none;
                padding: 0;
            }

            table {
                page-break-inside: auto;
            }

            tr {
                page-break-inside: avoid;
                page-break-after: auto;
            }

            th {
                position: static;
            }

            tr:hover {
                background: transparent;
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

	sb.WriteString(`                </tr>
            </thead>
            <tbody>
`)

	for _, item := range export.Items {
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
				sb.WriteString(html.EscapeString(*item.FolderID))
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

		sb.WriteString("                </tr>\n")
	}

	sb.WriteString(`            </tbody>
        </table>
    </div>
    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const table = document.querySelector('table');
            const headers = table.querySelectorAll('th');
            const tbody = table.querySelector('tbody');

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
