import { useEffect, useState } from "react";
import GamePage from "./components/GamePage";
import HomePage from "./components/HomePage";
import "./App.css";

function App() {
  const [isGamePage, setIsGamePage] = useState(() => window.location.pathname === "/play");
  useEffect(() => {
    const handlePopState = () => setIsGamePage(window.location.pathname === "/play");
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);
  const navigate = (path: "/" | "/play") => { window.history.pushState({}, "", path); setIsGamePage(path === "/play"); };
  return isGamePage ? <GamePage onHome={() => navigate("/")} /> : <HomePage onStart={() => navigate("/play")} />;
}

export default App;
