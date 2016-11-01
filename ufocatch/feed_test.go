package ufocatch

import (
	"encoding/xml"
	"fmt"
)

const content = `
<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>有報キャッチャー - EDINET情報配信サービス</title>
  <link href="http://resource.ufocatch.com/" />
  <updated>2016-10-08T06:41:41Z</updated>
  <id>http://resource.ufocatch.com/atom/edinetx/1</id>
  <link rel="self" href="http://resource.ufocatch.com/atom/edinetx/1" />
  <link rel="first" href="http://resource.ufocatch.com/atom/edinetx/1" />
  <link rel="last" href="http://resource.ufocatch.com/atom/edinetx/648" />
  <link rel="next" href="http://resource.ufocatch.com/atom/edinetx/2" />
  <entry>
    <title>【E20608】株式会社神明 変更報告書</title>
    <link rel="alternate" type="application/pdf" href="http://resource.ufocatch.com/pdf/edinet/ED2016100700001" />
    <id>ED2016100700001</id>
    <docid>S1008SOD</docid>
    <updated>2016-10-07T00:00:00+09:00</updated>
    <link rel="related" type="application/zip" href="http://resource.ufocatch.com/data/edinet/ED2016100700001" />
    <link rel="related" type="text/html" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/0000000_header_jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_ixbrl.htm" />
    <link rel="related" type="text/html" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/0101010_honbun_jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_ixbrl.htm" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07.xbrl" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07.xsd" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_def.xml" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_lab.xml" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_lab-en.xml" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/jplvh010000-lvh-001_E20608-000_2016-10-01_02_2016-10-07_pre.xml" />
    <link rel="related" type="text/xml" href="http://resource.ufocatch.com/xbrl/edinet/ED2016100700001/PublicDoc/manifest_PublicDoc.xml" />
  </entry>
</feed>
</xml>
	`

func ExampleFeed() {
	var feed Feed
	if err := xml.Unmarshal([]byte(content), &feed); err != nil {
		fmt.Println(err)
	}

	fmt.Println(feed.Updated)
	fmt.Println(len(feed.Links))
	fmt.Println(len(feed.Entries))
	// Output:
	// 2016-10-08 06:41:41 +0000 UTC
	// 5
	// 1
}

func ExampleFeedLink() {
	var feed Feed
	if err := xml.Unmarshal([]byte(content), &feed); err != nil {
		fmt.Println(err)
	}

	fmt.Println(feed.Link(""))
	fmt.Println(feed.Link("self"))
	fmt.Println(feed.Link("first"))
	fmt.Println(feed.Link("last"))
	fmt.Println(feed.Link("next"))
	fmt.Println(feed.Link("other"))
	// Output:
	// http://resource.ufocatch.com/
	// http://resource.ufocatch.com/atom/edinetx/1
	// http://resource.ufocatch.com/atom/edinetx/1
	// http://resource.ufocatch.com/atom/edinetx/648
	// http://resource.ufocatch.com/atom/edinetx/2
	// <nil>
}

func ExampleEntry() {
	var feed Feed
	if err := xml.Unmarshal([]byte(content), &feed); err != nil {
		fmt.Println(err)
	}
	link := feed.Entries[0].Links[0]
	fmt.Println(link.Rel)
	fmt.Println(link.Type)
	fmt.Println(link.URL)
	// Output:
	// alternate
	// application/pdf
	// http://resource.ufocatch.com/pdf/edinet/ED2016100700001
}
