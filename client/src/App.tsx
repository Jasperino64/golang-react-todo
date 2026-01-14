import { Container, Stack } from "@chakra-ui/react";
import Navbar from "./components/Navbar";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";

console.log(import.meta.env.MODE); // development
export const BASE_URL =
  import.meta.env.MODE === "development"
    ? "http://localhost:8080/api"
    : "https://golang-todo-api-production.up.railway.app/api";
console.log(BASE_URL); // http://localhost:8080/api

function App() {
  return (
    <>
      <Stack h="100vh">
        <Navbar />
        <Container>
          <TodoForm />
          <TodoList />
        </Container>
      </Stack>
    </>
  );
}

export default App;
