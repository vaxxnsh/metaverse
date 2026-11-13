type HomePageProps = { onStart: () => void; };

function HomePage({ onStart }: HomePageProps) {
  return <main className="home-page"><span className="eyebrow">A tiny Phaser 4 adventure</span><h1>Meadow Trail</h1><p>Take a short walk through a hand-built tile map.</p><button className="play-button" onClick={onStart} type="button">Start exploring <span>→</span></button><small>Keyboard: arrow keys or WASD</small></main>;
}

export default HomePage;
