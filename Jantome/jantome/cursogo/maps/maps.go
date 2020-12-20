package maps

//Devuelve Map
func GetMap() map[string]int {
	//miMap := make(map[string]int)

	//asignar valores directamente
	miMap := map[string]int{
		"Antonio": 30,
		"Maria":   31,
	}
	//se pueden serguir asignando valores independientemente asi
	miMap["edad1"] = 18
	miMap["edad2"] = 19

	//para eliminar algo del mapa
	delete(miMap, "edad1")

	return miMap

}

//le pasamos key y devuelve entero
func GetKeyMap(key string) int {

	//asignar valores directamente
	miMap := map[string]int{
		"Antonio": 30,
		"Maria":   31,
	}
	//se pueden serguir asignando valores independientemente asi
	miMap["edad1"] = 18
	miMap["edad2"] = 19

	//para eliminar algo del mapa
	delete(miMap, "edad1")

	return miMap[key]
}
