# DZone: Programming & DevOps Refcardz Downloader

## What is DZone.com

- DZone.com is one of the world's largest online communities and leading publisher of knowledge resources for software developers. Every day, hundreds of thousands of developers come to DZone.com to read about the latest technology trends and learn about new technologies, methodologies, and best practices through shared knowledge.

## Requirements

- Golang `1.10.2` or higher.
- [DZone](https://dzone.com) free account.

## Application

- Run `go run main.go`.

## Technical details of the solution

1. Loop through the assets list websites - `https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page=XX&sort=downloads&type=refcard`, where `XX` is `1` to `XX` (until this empty response is returned: `{"success":true,"result":{"data":{"assets":[],"sort":"downloads"}},"status":200}`). At the time of writing, there are `24` pages. See this example of the valid JSON response:

```json
{
  "success": true,
  "result": {
    "data": {
      "assets": [
        <---------- OMITTED ---------->
        {
          "id": 520107,
          "title": "GWT Style, Configuration and JSNI Reference",
          "details": "Introduces Ajax, a group interrelated techniques used in client-side web development for creating asynchronous web applications.",
          "subtitle": "Using the Google Web Toolkit",
          "collaborators": "Jill Tomich",
          "downloads": 29488,
          "views": 115116,
          "cover": "//dz2cdn3.dzone.com/storage/rc-covers/2806-dzone_refcard_.png",
          "host": null,
          "url": "/refcardz/gwt-style-configuration-and-js",
          "tags": [
            "frameworks",
            "javascript",
            "server-side",
            "java",
            "web dev",
            "ajax &amp; scripting"
          ],
          "color": "purple",
          "type": "refcard",
          "pdf": "/asset/download/6",
          "authors": [
            {
              "id": 327457,
              "name": "Robert Hansen",
              "avatar": "https://secure.gravatar.com/avatar/ae431e508cbc54620c27a0d612d4f93c?d=identicon&r=PG",
              "url": "/users/327457/rhansen1392.html"
            }
          ],
          "saveStatus": {
            "saved": false,
            "canSave": true,
            "count": 63
          }
        }
        <---------- OMITTED ---------->
      ],
      "sort": "downloads"
    }
  },
  "status": 200
}
```

2. On each of these asset list websites, extract the following information from the returned JSON file:

    - Title:      `'.result.data.assets[].title'`.
    - PDF suffix: `'.result.data.assets[].pdf'`.

3. Download the PDF by prefixing the PDF link with `https://dzone.com`, creating e.g. `https://dzone.com/asset/download/279342` and save it as `<Title>.pdf`, e.g. `GWT_Style,_Configuration_and_JSNI_Reference.pdf`

## Stats

|Item          |Size      |
|--------------|----------|
|Refcardz      |284       |
|All files size|000 B (GB)|
