package contacts

const Email = "support@example.com"

var support string 

func SetSupport(s string) { 
    support = s
}

func GetContact() string {
    return fmt.Sprintf("%s <%s>", support, Email)
} 