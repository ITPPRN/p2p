package utils


///////// auth ///////////
func ConvertInterfaceSliceToStringSlice(slice []interface{}) []string {
    strSlice := make([]string, len(slice))
    for i, v := range slice {
        // ต้องมั่นใจว่า v สามารถแปลงเป็น string ได้
        strSlice[i] = v.(string)
    }
    return strSlice
}

func GetSafeString(claims map[string]interface{}, key string) string {
    if v, ok := claims[key]; ok && v != nil {
        if s, ok := v.(string); ok {
            return s
        }
    }
    return ""
}