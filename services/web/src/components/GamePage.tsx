import GameCanvas from "./GameCanvas";

type GamePageProps = { onHome: () => void; };

function GamePage({ onHome }: GamePageProps) {
  return <main className="game-page"><button className="back" onClick={onHome} type="button">← Home</button><GameCanvas /><p>Move with arrow keys or WASD. Avoid the river and boulders.</p></main>;
}

export default GamePage;
