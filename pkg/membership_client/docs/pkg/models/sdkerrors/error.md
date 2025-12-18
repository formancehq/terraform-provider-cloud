# Error


## Fields

| Field                                                   | Type                                                    | Required                                                | Description                                             |
| ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| `ErrorCode`                                             | *string*                                                | :heavy_check_mark:                                      | N/A                                                     |
| `ErrorMessage`                                          | **string*                                               | :heavy_minus_sign:                                      | N/A                                                     |
| `RawResponse`                                           | [*http.Response](https://pkg.go.dev/net/http#Response)  | :heavy_minus_sign:                                      | Raw HTTP response; suitable for custom response parsing |