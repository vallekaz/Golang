package main 
import "fmt" 
const bienvenido string = "Bienvenido a la calculadora del Curso de Go" 
const msg string = "¿Qué deseas hacer?" 
const errorOpcion string = "Error: Opción inválida" 
const resultado string = "El resultado de %s los valores introducidos es: %d" 
const sumarText = "sumar" const restarText = "restar" 
const multiplicarText = "multiplicar" 
const dividirText = "dividir" 
func main() { 
	opcion := 0 fmt.Println("") 
	fmt.Println(bienvenido) 
	fmt.Println("") 
	fmt.Println(msg) 
	fmt.Println("") 
	fmt.Println("1. Sumar") 
	fmt.Println("2. Restar") 
	fmt.Println("3. Multiplicar") 
	fmt.Println("4. Dividir") 
	fmt.Println("") 
	fmt.Scanf("%d", &opcion) 
	// Usamos %d para capturar un número en lugar de %s que es para string 
	switch opcion { 
		case 1: 
		fmt.Printf(resultado, sumarText, suma()) 
		break 
		case 2: 
		fmt.Printf(resultado, restarText, resta()) 
		break 
		case 3: 
		fmt.Printf(resultado, multiplicarText, multiplicar()) 
		break 
		case 4: 
		fmt.Printf(resultado, dividirText, dividir()) 
		break 
		default: 
		fmt.Println(errorOpcion) 
		fmt.Println("") 
		} 

func capturarValores() (int, int) { 
	var ( a = 0 b = 0 ) 
	fmt.Print("Introduce el primer número: ") 
	fmt.Scanf("%d", &a) 
	fmt.Print("Introduce el segundo número: ") 
	fmt.Scanf("%d", &b) 
	return a, b 
	} 

func suma() int { 
	a,b := capturarValores() 
	return a + b 
} 

func resta() int { 
	a,b := capturarValores() 
	return a - b 
} 
func multiplicar() int { 
	a,b := capturarValores() 
	return a * b 
} 
func dividir() int { 
	a,b := capturarValores() 
	return a / b 
}