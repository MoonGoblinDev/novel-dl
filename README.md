# novel-dl

Browser-based novel downloader for 📖🐰 (stable)

## Features

- Novel download from 📖🐰 Booktoons (북토끼)
- Automatic file merging
- Quick execution via browser console or bookmark
- No program installation required

## Usage (Please star the repository above!)

### Via Bookmark (Recommended for repeated use)
1. Copy the script below (use the copy button on the right):
    ```javascript
    javascript:(function(){fetch('https://raw.githubusercontent.com/MoonGoblinDev/novel-dl/refs/heads/main/script.js').then(response=>{if(!response.ok){throw new Error(`Failed to fetch script: ${response.statusText}`);}return response.text();}).then(scriptContent=>{const script=document.createElement('script');script.textContent=scriptContent;document.head.appendChild(script);console.log('Script loaded and executed.');}).catch(error=>{console.error(error);});})();
    ```
2. In your browser, display the bookmarks bar by pressing `ctrl+shift+b`.
3. Add a bookmark to any page by pressing `ctrl+d`.
4. Right-click the bookmark and select "Edit".
5. Paste the copied script into the "URL" field of the bookmark (do *not* change the title).
6. Click the bookmark on the chapter list page of the novel you want to download.

### Via Browser Console (For one-time use, not recommended)
1. Copy [script.js](https://raw.githubusercontent.com/yeorinhieut/novel-dl/main/script.js).
2. On the chapter list page of the novel you want to download, open the browser console by pressing `f12` or `ctrl+shift+i`.
3. Paste the copied script into the console.

## FAQ

### Can I download multiple novels at once?

Trying to download multiple novels from different tabs simultaneously might trigger website blocking policies and prevent access. This is not recommended. (More than 60 requests per minute may be blocked.)

### I'm getting an error.

Please report the error in the [issues](https://github.com/yeorinhieut/novel-dl/issues) section or contribute by creating a pull request.

### I have a suggestion for improvement.

Please submit your suggestion in the [issues](https://github.com/yeorinhieut/novel-dl/issues) section or contribute by creating a pull request.

---
## Commissioning Other Programs
- We welcome inquiries for programs with public interest.
- We also offer custom software development for personal use.
- Please contact us at [email](mailto:yeorinhieut@gmail.com) or Discord yeorinhieut.
