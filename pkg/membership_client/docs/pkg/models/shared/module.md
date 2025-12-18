# Module


## Fields

| Field                                                                | Type                                                                 | Required                                                             | Description                                                          |
| -------------------------------------------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `Name`                                                               | *string*                                                             | :heavy_check_mark:                                                   | N/A                                                                  |
| `State`                                                              | [shared.ModuleState](../../../pkg/models/shared/modulestate.md)      | :heavy_check_mark:                                                   | N/A                                                                  |
| `Status`                                                             | [shared.ModuleStatus](../../../pkg/models/shared/modulestatus.md)    | :heavy_check_mark:                                                   | N/A                                                                  |
| `LastStatusUpdate`                                                   | [time.Time](https://pkg.go.dev/time#Time)                            | :heavy_check_mark:                                                   | N/A                                                                  |
| `LastStateUpdate`                                                    | [time.Time](https://pkg.go.dev/time#Time)                            | :heavy_check_mark:                                                   | N/A                                                                  |
| `ClusterStatus`                                                      | [*shared.ClusterStatus](../../../pkg/models/shared/clusterstatus.md) | :heavy_minus_sign:                                                   | N/A                                                                  |