package main

/*
Desarrollado por:
-> Barros Martínez Luis Enrique
-> Baustista Rios Alfredo
-> Cortes Velazquez Samuel Alejandro
──────▄▀▄─────▄▀▄
─────▄█░░▀▀▀▀▀░░█▄
─▄▄──█░░░░░░░░░░░█──▄▄
█▄▄█─█░░▀░░┬░░▀░░█─█▄▄█
Para la materia de Compiladores (5CM4) impartida por:
Gabriel de Jesus Rodriguez Jordan
╭━━━━━━━━━╮╭┓┈┈┏╮
┃┈┈┈┈┈┈┈┈┈┃┃╰╮╭╯┃
┃┈┈┈┈┈┈┈┈┈┃╰━┓┏━╯
┃┈┈┈┈┈0┈┈┈╰━━╯┃
┣━━━┳┳┳╮┈┈┈┈┈┈┃
╰━━━╰━╯━━━━━━━╯
*/
import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Nodo struct {
	hijos    map[byte]*Nodo
	completo bool
	valor    []byte
}

const (
	Terminal     = 0b1000_0000
	T_select     = 0b1100_0000
	T_distinct   = 0b1100_0001
	T_from       = 0b1000_0010
	T_id         = 0b1110_0011
	T_id_epsilon = 0b1111_0011
	T_dot        = 0b1100_0100
	T_comma      = 0b1100_0101
	T_star       = 0b1110_0110
	T_End        = 0b1111_1111 //255
)
const (
	No_Terminal = 0b0000_0000
	NT_Q        = 0b0000_0000
	NT_D        = 0b0000_0001
	NT_P        = 0b0000_0010
	NT_A        = 0b0000_0011
	NT_A1       = 0b0000_0100
	NT_A2       = 0b0000_0101
	NT_A3       = 0b0000_0110
	NT_T        = 0b0000_0111
	NT_T1       = 0b0000_1000
	NT_T2       = 0b0000_1001
	NT_T3       = 0b0000_1010
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Error: ", p)
		}
	}()
	Args := os.Args
	Args = Args[1:]
	if len(Args) < 1 {
		fmt.Println("Error, No file!!")
	} else {

		byte_data, err := os.ReadFile(Args[0])
		if err != nil {
			log.Fatal(err)
		}
		string_data := string(byte_data)
		tokens := lexate(string_data)
		//fmt.Println(tokens)
		//Analizador  Desendente
		fmt.Println(ASD(tokens))
		//Analizador Ascendente
		fmt.Println(ASA(tokens))

	} //hola
}
func ASD(tokens_file string) bool {
	//Remplazar los tokens por valores numericos
	tokens := strings.Split(tokens_file, " ") //separando tokens por espacios
	tkns := make([]byte, 0)                   //arreglo de bites, se guardan los tokens en su valor binario
	for _, token := range tokens {
		switch token {
		case "select":
			tkns = append(tkns, T_select)
		case "distinct":
			tkns = append(tkns, T_distinct)
		case "from":
			tkns = append(tkns, T_from)
		case "id":
			tkns = append(tkns, T_id)
		case ".":
			tkns = append(tkns, T_dot)
		case ",":
			tkns = append(tkns, T_comma)
		case "*":
			tkns = append(tkns, T_star)
		default:
			fmt.Println([]rune(token))
			panic("Token no reconocido: " + token)
		}
	}
	pila := make([]byte, 0)    // TOPE = v
	pila = append(pila, NT_Q)  // Pila = [NT_Q]
	pila = append(pila, T_End) // Pila = [NT_Q, T_End]
	tam_entrada := len(tkns)
	tkn_actual := 0
	seguro := 0
	//fmt.Println(tkns, "\n----------")
	for tkn_actual < tam_entrada { //Mientras no sea el ultimo tokn
		//fmt.Println("Pila:", pila, "\nTokens:", tkns[tkn_actual:], "\n-------------------")
		if pila[0] < Terminal { //Si el tope de la pila no es un terminal
			seguro = 0
			for pila[0] < Terminal { // Mientras el tope de la pila no sea un terminal
				pila = append(Expandir(pila[0], tkns[tkn_actual]), pila[1:]...) //expande la pila
			}
		} else if pila[0] == tkns[tkn_actual] {
			seguro = 0
			pila = pila[1:] //Cortar tope de la pila
			tkn_actual++    // Recorrer entrada
		} else {
			seguro++
			if seguro >= tam_entrada {
				return false
			}
		}
		if pila[0] == T_End { //Tope de la pila igual a terminal, salimos del bucle
			break
		}
	}
	//fmt.Println(pila)
	for pila[0] < Terminal { //Si el tope de la pila no es un terminal
		pila = append(Expandir(pila[0], 0), pila[1:]...) //Reemplaza por epsilon
		//fmt.Println(pila)
	}
	//fmt.Println(pila)
	if pila[0] == T_End && tkn_actual >= len(tkns) { // verifica si la pila esta vacia y que el token actual llego al final de la lista de tokens
		return true //Analisis correcto carajo
	} else {
		return false //esta mal escrito papito, asi de huevos? como?
	}
}

func Expandir(token byte, entrada byte) []byte { //toma el token NT a expandir y el ultimo T que se ha leido
	//	fmt.Println(token)
	switch token {
	case NT_Q:
		return []byte{T_select, NT_D, T_from, NT_T}
	case NT_D:
		switch entrada {
		case T_distinct:
			return []byte{T_distinct, NT_P}
		default:
			return []byte{NT_P}
		}
	case NT_P:
		switch entrada {
		case T_star:
			return []byte{T_star}
		default:
			return []byte{NT_A}
		}
	case NT_A:
		return []byte{NT_A2, NT_A1}
	case NT_A1:
		switch entrada {
		case T_comma:
			return []byte{T_comma, NT_A}
		default:
			return []byte{} //epsilon
		}
	case NT_A2:
		return []byte{T_id, NT_A3}
	case NT_A3:
		switch entrada {
		case T_dot:
			return []byte{T_dot, T_id}
		default:
			return []byte{} //epsilon
		}
	case NT_T:
		return []byte{NT_T2, NT_T1}
	case NT_T1:
		switch entrada {
		case T_comma:
			return []byte{T_comma, NT_T}
		default:
			return []byte{} //epsilon
		}
	case NT_T2:
		return []byte{T_id, NT_T3}
	case NT_T3:
		switch entrada {
		case T_id:
			return []byte{T_id}
		default:
			return []byte{} //epsilon
		}

	}
	panic("Salto no efectuado o no contemplado => tkn:" + strconv.Itoa(int(token)) + " entrada:" + strconv.Itoa(int(entrada)))
	// no le sabes papito, escribe bien otra vez
}

func ASA(tokens_file string) bool {
	//Remplazar los tokens por valores numericos
	tokens := strings.Split(tokens_file, " ")
	tkns := make([]byte, 0)
	for _, token := range tokens {
		switch token {
		case "select":
			tkns = append(tkns, T_select)
		case "distinct":
			tkns = append(tkns, T_distinct)
		case "from":
			tkns = append(tkns, T_from)
		case "id":
			tkns = append(tkns, T_id)
		case ".":
			tkns = append(tkns, T_dot)
		case ",":
			tkns = append(tkns, T_comma)
		case "*":
			tkns = append(tkns, T_star)
		default:
			fmt.Println([]rune(token))
			panic("Token no reconocido: " + token)
		}
	}
	reducciones := arbol_reducciones()
	pila := make([]byte, 0)
	//pila = append(pila, T_End)
	tam_tokns := len(tkns)
	tkn_actual := 0
	seguro := 0
	post_from := false
	for tkn_actual < tam_tokns {
		//fmt.Println("pila,tkns:", pila, tkns[tkn_actual:])
		switch tkns[tkn_actual] {
		case T_select:
			Desplazar(&pila, &tkns, &tkn_actual)
			seguro = 0
		case T_distinct:
			Desplazar(&pila, &tkns, &tkn_actual)
		case T_from:
			if pila[len(pila)-1] == NT_D {
				Desplazar(&pila, &tkns, &tkn_actual)
				post_from = true

				seguro = 0
			} else {
				reducciones.Reducir(&pila, tkns[tkn_actual], 0)
				seguro++
				if seguro >= tam_tokns {
					return false
				}
			}
		case T_comma:
			if !post_from && pila[len(pila)-1] != T_id && pila[len(pila)-1] != NT_A1 {
				Desplazar(&pila, &tkns, &tkn_actual)
				seguro = 0
			} else if post_from && pila[len(pila)-1] != T_id_epsilon && pila[len(pila)-1] != NT_T1 {
				Desplazar(&pila, &tkns, &tkn_actual)
				//pila[len(pila)-1] += 0b0001_0000
				seguro = 0
			} else {
				reducciones.Reducir(&pila, tkns[tkn_actual], 0)
				seguro++
				if seguro >= tam_tokns {
					return false
				}
			}
		case T_id:
			//fmt.Println(pila, tkns[tkn_actual:], ",", tkns[tkn_actual])
			Desplazar(&pila, &tkns, &tkn_actual)
			if post_from {
				pila[len(pila)-1] += 0b0001_0000
			}
			seguro = 0
		case T_dot:
			if pila[len(pila)-1] == T_id {
				Desplazar(&pila, &tkns, &tkn_actual)
				seguro = 0
			} else {
				reducciones.Reducir(&pila, tkns[tkn_actual], 0)
				seguro++
				if seguro >= tam_tokns {
					return false
				}
			}
		case T_star:
			Desplazar(&pila, &tkns, &tkn_actual)

		default:
			reducciones.Reducir(&pila, tkns[tkn_actual], 0)
			seguro++
			if seguro >= tam_tokns {
				return false
			}
		}

	}
	seguro = 0
	for reducciones.Reducir(&pila, T_End, 0) {
		seguro++
		if seguro >= tam_tokns {
			return false
		}
	}
	//fmt.Println("pila,tkns:", pila, tkns[tkn_actual:])

	return pila[0] == NT_Q

}
func Desplazar(pila *[]byte, tkns *[]byte, tkn_actual *int) {
	*pila = append(*pila, (*tkns)[*tkn_actual])
	*tkn_actual++

}

func (Reducciones *Nodo) Reducir(pila *[]byte, entrada byte, profundidad uint8) bool {
	// Obtener la reduccin utilizando la funcin buscar_reduccion
	if len(*pila) == 0 {
		return false
	}
	if profundidad < 4 {
		corte := 4
		if len(*pila) < 4 {
			corte = len(*pila)
		} else {
			corte -= int(profundidad)
		}
		//fmt.Println("Longitud pila de entrada: ", len(*pila))
		Reduccion, longitud := Reducciones.buscar_reduccion((*pila)[len(*pila)-corte:])
		//fmt.Println("Reduccion cachada : ", Reduccion, "longitud:", longitud)
		if longitud != 0 {
			//fmt.Println("d1")
			(*pila) = append((*pila)[:len((*pila))-longitud], Reduccion...)
			//fmt.Println("pila:", pila)
			return true
		}
		//}
	}
	//fmt.Println("falsee")
	return false // No se pudo realizar una reduccin
}

func (raiz *Nodo) buscar_reduccion(pila []byte) ([]byte, int) {

	cantidad_tkns := 4
	if len(pila) < 4 {
		cantidad_tkns = len(pila)
	}
	//fmt.Println("cantidad tkns:", cantidad_tkns)
	for i := cantidad_tkns; i > 0; i-- {
		reduccion, longitud := raiz.buscarReduccionRecursiva(pila[cantidad_tkns-i:], 0)
		//fmt.Println("reduccion encontrada: ", reduccion, longitud)
		if longitud != 0 {

			return reduccion, longitud
		}
	}

	return nil, 0
}

func (nodoActual *Nodo) buscarReduccionRecursiva(pila []byte, profundidad uint8) ([]byte, int) {
	//fmt.Println("Pila de busqueda:", pila)
	if len(pila) == 0 {
		if nodoActual.completo {
			return nodoActual.valor, int(profundidad)
		}
		return nil, 0
	}
	token := pila[len(pila)-1]
	hijo, encontrado := nodoActual.hijos[token]
	if !encontrado {
		return nil, 0
	}
	penultimo := len(pila) - 1
	return hijo.buscarReduccionRecursiva(pila[:penultimo], profundidad+1)
}

func arbol_reducciones() *Nodo {
	Root := Nodo{
		valor:    nil,
		completo: false,
		hijos: map[byte]*Nodo{
			//-------------------------------------------- Transicion id ->Nodo A1
			T_id: {
				valor:    []byte{NT_A1},
				completo: true,
				hijos: map[byte]*Nodo{
					//------------------------------------ Transicion dot -> Nodo A2
					T_dot: {
						valor:    []byte{NT_A2},
						completo: true,
						hijos: map[byte]*Nodo{
							//---------------------------- Transicion id -> Nodo A1
							T_id: {
								valor:    []byte{NT_A1},
								completo: true,
								hijos:    map[byte]*Nodo{},
							},
						},
					},
					T_id: {
						valor:    []byte{NT_T1},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
				},
			},
			//-------------------------------------------- Transicion idE -> Nodo T1
			T_id_epsilon: {
				valor:    []byte{NT_T1},
				completo: true,
				hijos: map[byte]*Nodo{
					//------------------------------------ Transicion id -> Nodo T1
					T_id: {
						valor:    []byte{NT_T1},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
					T_id_epsilon: {
						valor:    []byte{NT_T1},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
				},
			},
			//-------------------------------------------- Transicion T1 -> Nodo T
			NT_T1: {
				valor:    []byte{NT_T},
				completo: true,
				hijos: map[byte]*Nodo{
					//------------------------------------ Transicion coma -> Nodo vacio
					T_comma: {
						completo: false,
						valor:    []byte{},
						hijos: map[byte]*Nodo{
							//---------------------------- Transicion T1 -> T
							NT_T1: {
								completo: true,
								valor:    []byte{NT_T},
								hijos:    map[byte]*Nodo{},
							},
							NT_T: {
								completo: true,
								valor:    []byte{NT_T},
								hijos:    map[byte]*Nodo{},
							},
						},
					},
				},
			},
			//--------------------------------------------Transicion T2 -> Nodo void
			NT_T2: {
				valor:    []byte{},
				completo: false,
				hijos: map[byte]*Nodo{
					//------------------------------------Transicion T1 -> Nodo T1?
					T_id: {
						valor:    []byte{T_id},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
				},
			},
			//--------------------------------------------Transicion A1 -> Nodo A
			NT_A1: {
				valor:    []byte{NT_A},
				completo: true,
				hijos: map[byte]*Nodo{
					//------------------------------------Transicion comma -> Nodo void
					T_comma: {
						valor:    []byte{},
						completo: false,
						hijos: map[byte]*Nodo{
							//----------------------------Transicion A -> Nodo A
							NT_A: {
								valor:    []byte{NT_A},
								completo: true,
								hijos:    map[byte]*Nodo{},
							},
						},
					},
				},
			},
			//--------------------------------------------Transicion A2 -> Nodo void
			NT_A2: {
				valor:    []byte{},
				completo: false,
				hijos: map[byte]*Nodo{
					//-----------------------------------Transicion id -> A1
					T_id: {
						valor:    []byte{NT_A1},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
				},
			},
			//-------------------------------------------Transicion A -> Nodo P
			NT_A: {
				valor:    []byte{NT_P},
				completo: true,
				hijos:    map[byte]*Nodo{},
			},
			//------------------------------------------Transicion P -> Nodo D
			NT_P: {
				valor:    []byte{NT_D},
				completo: true,
				hijos: map[byte]*Nodo{
					//----------------------------------Transicion distinct -> Nodo D
					T_distinct: {
						valor:    []byte{NT_D},
						completo: true,
						hijos:    map[byte]*Nodo{},
					},
				},
			},
			//-------------------------------------------Transicion T -> Nodo void
			NT_T: {
				valor:    []byte{},
				completo: false,
				hijos: map[byte]*Nodo{
					//-----------------------------------Transicion from -> Nodo void
					T_from: {
						valor:    []byte{},
						completo: false,
						hijos: map[byte]*Nodo{
							//---------------------------Transicion D -> Nodo void
							NT_D: {
								valor:    []byte{},
								completo: false,
								hijos: map[byte]*Nodo{
									//-------------------Transicion select -> Nodo Q
									T_select: {
										valor:    []byte{NT_Q},
										completo: true,
										hijos:    map[byte]*Nodo{},
									},
								},
							},
						},
					},
				},
			},
			T_star: {
				valor:    []byte{NT_P},
				completo: true,
				hijos:    map[byte]*Nodo{},
			},
		},
	}

	return &Root
}

func lexate(sqlCode string) string {
	// Definir expresiones regulares para cada token
	selectRegex := regexp.MustCompile(`\bselect\b`)
	fromRegex := regexp.MustCompile(`\bfrom\b`)
	dotRegex := regexp.MustCompile(`\.`)
	commaRegex := regexp.MustCompile(`,`)
	distinctRegex := regexp.MustCompile(`\bdistinct\b`)
	starRegex := regexp.MustCompile(`\*`)
	
	// Reemplazar ", " por " , "
	sqlCode = strings.ReplaceAll(sqlCode, ", ", " , ")
	// Reemplazar ". " por " . "
	sqlCode = strings.ReplaceAll(sqlCode, ". ", " . ")
	// Dividir el código en palabras
	words := strings.Fields(sqlCode)

	// Iterar sobre las palabras y clasificarlas en tokens
	var resultTokens []string
	for _, word := range words {
		switch {
		case selectRegex.MatchString(word):
			resultTokens = append(resultTokens, "select")
		case fromRegex.MatchString(word):
			resultTokens = append(resultTokens, "from")
		case dotRegex.MatchString(word):
			resultTokens = append(resultTokens, ".")
		case commaRegex.MatchString(word):
			resultTokens = append(resultTokens, ",")
		case distinctRegex.MatchString(word):
			resultTokens = append(resultTokens, "distinct")
		case starRegex.MatchString(word):
			resultTokens = append(resultTokens, "*")
		default:
			resultTokens = append(resultTokens, "id")
		}
	}

	// Unir los tokens en un string separado por espacios
	result := strings.Join(resultTokens, " ")
	return result
}
