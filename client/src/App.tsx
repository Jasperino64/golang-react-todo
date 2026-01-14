import { Container, Stack } from "@chakra-ui/react";
import Navbar from "./components/Navbar";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";
import { API_BASE_URL, apiFetch } from "@/api/client";
import { useEffect } from "react";

console.log(import.meta.env.MODE); // development
export const BASE_URL = API_BASE_URL;
console.log(BASE_URL); // http://localhost:8080/api

function App() {
  useEffect(() => {
    apiFetch("/health")
      .then((r: Response) => r.text())
      .then(console.log);
  }, []);

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
