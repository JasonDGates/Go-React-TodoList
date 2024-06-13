import { useState } from "react";
import { Button, Stack } from "@chakra-ui/react";

function App() {
  const [count, setCount] = useState(0);

  return (
    <Stack h="100vh">
      <Button>Hello</Button>
    </Stack>
  );
}

export default App;
