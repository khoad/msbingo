# msbingo
An implementation of NBFX and NBFS, which is expressed in HTTP terms as `Content-Type: application/soap+msbin1`, written in pure Go to enable interop with a WCF service and Go with no other dependencies (like Mono or Windows).

This implementation follows the Microsoft specification closely in naming, structure and example. Tests for both decoding and encoding have been written for each of the Structure Examples provided, to validate the given bytes encode to the given XML and vice versa.

Currently the decoding side is more (nearly!) complete in terms of implementation of the individual record types. You can see the current state by running the tests.

Encoding is less complete, as it is not necessary to encode exactly as .NET WCF would. It is only necessary to encode properly, such that the target service can decode the original XML message properly.

# Usage

``` go
url := fmt.Sprintf("%s/Path/To/ExampleService.svc", s.apiBaseUrl)

xmlInput := "<s:Envelope xmlns:a=\"http://www.w3.org/2005/08/addressing\" xmlns:s=\"http://www.w3.org/2003/05/soap-envelope\"><s:Header><a:Action s:mustUnderstand=\"1\">action</a:Action></s:Header><s:Body><Inventory>0</Inventory></s:Body></s:Envelope>"

encodedXml, err := nbfs.NewEncoder().Encode(bytes.NewBufferString(xmlInput))
if err != nil {
	// handle encoding error
}

req, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedXml))
if err != nil {
	// handle request creation error
}

req.Header.Add("Content-Type", "application/soap+msbin1")
resp, err := httpClient.Do(req)
if err != nil {
	// handle response error
}
defer resp.Body.Close()

xmlRes, err := nbfs.NewDecoder().Decode(resp.Body)
if err != nil {
	// handle decoding error
}

// do something with your decoded xml response
```

# Codec details
NBFS is a codec developed by Microsoft for use primarily by WCF webservices. It is essentially a binary encoding for Soap XML messages optimized for reducing bytes over the wire.  They have published the specification in multiple parts:
* [NBFX (.NET Binary Format: Xml)](https://msdn.microsoft.com/en-us/library/cc219210.aspx)
* [NBFS (.NET Binary Format: SOAP Data Structure)](https://msdn.microsoft.com/en-us/library/cc219175.aspx)
where NBFS is essentially NBFX with standard DictionaryString entries for strings commonly used in SOAP, such as "Envelope", "http://www.w3.org/2003/05/soap-envelope/", etc., to minimize the bytewise size overhead of the SOAP protocol.
