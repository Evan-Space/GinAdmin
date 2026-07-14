import { createRoot } from "react-dom/client";
import App from "./App";
import "@src/style/global.css";

createRoot(document.getElementById("root") as HTMLElement).render(<App />);
