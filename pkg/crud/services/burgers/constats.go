package burgers

const updateInTable = "UPDATE burgers SET removed = true where id = $1;"
const insertToTable = "INSERT INTO burgers(name, price) VALUES ($1, $2);"
const searchWearFalse = "SELECT id, name, price FROM burgers WHERE removed = FALSE"