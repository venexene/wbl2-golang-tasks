// Package wget provides functionality for mirroring web pages
package wget

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func setAttr(node *html.Node, key string, value string) {
	for i, attr := range node.Attr {
		if attr.Key == key {
			node.Attr[i].Val = value
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
}

func resolveURL(href string, base *url.URL) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return base.ResolveReference(uri).String()
}

func generateFilename(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Sprintf("resource_%x", md5.Sum([]byte(urlStr)))
	}

	name := path.Base(parsedURL.Path)
	if name == "." || name == "/" || name == "" {
		name = "index"
	}

	if strings.Contains(name, "?") {
		name = strings.Split(name, "?")[0]
	}

	ext := path.Ext(name)
	if ext == "" {
		switch {
		case strings.Contains(parsedURL.Path, ".css") || strings.Contains(urlStr, ".css"):
			name += ".css"
		case strings.Contains(parsedURL.Path, ".js") || strings.Contains(urlStr, ".js"):
			name += ".js"
		case strings.Contains(parsedURL.Path, ".png") || strings.Contains(urlStr, ".png"):
			name += ".png"
		case strings.Contains(parsedURL.Path, ".jpg") || strings.Contains(urlStr, ".jpg"):
			name += ".jpg"
		case strings.Contains(parsedURL.Path, ".svg") || strings.Contains(urlStr, ".svg"):
			name += ".svg"
		default:
			name += ".bin"
		}
	}

	return name
}

func getResourceLinks(node *html.Node, baseURL *url.URL) map[string]string {
	links := make(map[string]string)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var resourceURL string

			switch n.Data {
			case "link":
				if getAttr(n, "rel") == "stylesheet" {
					resourceURL = getAttr(n, "href")
				}
			case "script":
				resourceURL = getAttr(n, "src")
			case "img":
				resourceURL = getAttr(n, "src")
			case "source":
				resourceURL = getAttr(n, "src")
			case "audio":
				resourceURL = getAttr(n, "src")
			case "video":
				resourceURL = getAttr(n, "src")
			case "embed":
				resourceURL = getAttr(n, "src")
			case "object":
				resourceURL = getAttr(n, "data")
			}

			if resourceURL != "" {
				absoluteURL := resolveURL(resourceURL, baseURL)
				if absoluteURL != "" {
					filename := generateFilename(absoluteURL)
					links[absoluteURL] = filename
				}
			}

			if n.Data == "img" || n.Data == "source" {
				if srcset := getAttr(n, "srcset"); srcset != "" {
					entries := strings.Split(srcset, ",")
					for _, entry := range entries {
						parts := strings.Split(strings.TrimSpace(entry), " ")
						if len(parts) > 0 && parts[0] != "" {
							absoluteURL := resolveURL(parts[0], baseURL)
							if absoluteURL != "" {
								filename := generateFilename(absoluteURL)
								links[absoluteURL] = filename
							}
						}
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)
	return links
}

func downloadResource(resURL, filePath string) error {
	response, err := http.Get(resURL)
	if err != nil {
		return fmt.Errorf("Error downloading resource: %w", err)
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating resource file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("Error writing resource file: %w", err)
	}

	return nil
}

func replaceResourceLinks(node *html.Node, resLinks map[string]string, baseURL *url.URL) {
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "link":
				if getAttr(n, "rel") == "stylesheet" {
					if href := getAttr(n, "href"); href != "" {
						absoluteURL := resolveURL(href, baseURL)
						if filename, exists := resLinks[absoluteURL]; exists {
							setAttr(n, "href", filepath.Join("resources", filename))
						}
					}
				}

			case "script":
				if src := getAttr(n, "src"); src != "" {
					absoluteURL := resolveURL(src, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "src", filepath.Join("resources", filename))
					}
				}

			case "img":
				if src := getAttr(n, "src"); src != "" {
					absoluteURL := resolveURL(src, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "src", filepath.Join("resources", filename))
					}
				}
				if srcset := getAttr(n, "srcset"); srcset != "" {
					entries := strings.Split(srcset, ",")
					var newEntries []string
					for _, entry := range entries {
						parts := strings.Split(strings.TrimSpace(entry), " ")
						if len(parts) > 0 && parts[0] != "" {
							absoluteURL := resolveURL(parts[0], baseURL)
							if filename, exists := resLinks[absoluteURL]; exists {
								newPath := filepath.Join("resources", filename)
								if len(parts) > 1 {
									newEntries = append(newEntries, newPath+" "+parts[1])
								} else {
									newEntries = append(newEntries, newPath)
								}
							} else {
								newEntries = append(newEntries, entry)
							}
						}
					}
					if len(newEntries) > 0 {
						setAttr(n, "srcset", strings.Join(newEntries, ", "))
					}
				}

			case "source":
				if src := getAttr(n, "src"); src != "" {
					absoluteURL := resolveURL(src, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "src", filepath.Join("resources", filename))
					}
				}
				if srcset := getAttr(n, "srcset"); srcset != "" {
					entries := strings.Split(srcset, ",")
					var newEntries []string
					for _, entry := range entries {
						parts := strings.Split(strings.TrimSpace(entry), " ")
						if len(parts) > 0 && parts[0] != "" {
							absoluteURL := resolveURL(parts[0], baseURL)
							if filename, exists := resLinks[absoluteURL]; exists {
								newPath := filepath.Join("resources", filename)
								if len(parts) > 1 {
									newEntries = append(newEntries, newPath+" "+parts[1])
								} else {
									newEntries = append(newEntries, newPath)
								}
							} else {
								newEntries = append(newEntries, entry)
							}
						}
					}
					if len(newEntries) > 0 {
						setAttr(n, "srcset", strings.Join(newEntries, ", "))
					}
				}

			case "audio", "video":
				if src := getAttr(n, "src"); src != "" {
					absoluteURL := resolveURL(src, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "src", filepath.Join("resources", filename))
					}
				}

			case "embed":
				if src := getAttr(n, "src"); src != "" {
					absoluteURL := resolveURL(src, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "src", filepath.Join("resources", filename))
					}
				}

			case "object":
				if data := getAttr(n, "data"); data != "" {
					absoluteURL := resolveURL(data, baseURL)
					if filename, exists := resLinks[absoluteURL]; exists {
						setAttr(n, "data", filepath.Join("resources", filename))
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)
}

// MirrorPage creates local copy of web page
func MirrorPage(link string) error {
	response, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("Error getting response by link: %w", err)
	}

	parsedLink, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("Error parsing link: %w", err)
	}

	folderName := parsedLink.Hostname()
	err = os.Mkdir(folderName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating main folder: %w", err)
	}

	resFolder := filepath.Join(folderName, "resources")
	err = os.Mkdir(resFolder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating resources folder: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading body: %w", err)
	}

	parsedHTML, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("Error parsing HTML: %w", err)
	}

	resLinks := getResourceLinks(parsedHTML, parsedLink)
	fmt.Printf("Found %d resources:\n", len(resLinks))

	for resURL, filename := range resLinks {
		fmt.Printf("Downloading %s -> %s\n", resURL, filename)
		filePath := filepath.Join(resFolder, filename)
		err := downloadResource(resURL, filePath)
		if err != nil {
			fmt.Printf("!Warning: Failed to download %s: %v\n", resURL, err)
			delete(resLinks, resURL)
		}
	}

	replaceResourceLinks(parsedHTML, resLinks, parsedLink)

	file, err := os.Create(filepath.Join(folderName, "index.html"))
	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}
	defer file.Close()

	err = html.Render(file, parsedHTML)
	if err != nil {
		return fmt.Errorf("Error writing HTML: %w", err)
	}

	fmt.Printf("Successfully mirrored page to %s/\n", folderName)
	return nil
}
