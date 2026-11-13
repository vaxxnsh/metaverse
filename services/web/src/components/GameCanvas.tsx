import { useEffect, useRef } from "react";
import * as Phaser from "phaser";
import MeadowScene from "../game/MeadowScene";
import { GAME_HEIGHT, GAME_WIDTH } from "../lib/game-constants";

function GameCanvas() {
  const host = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!host.current) return;
    const game = new Phaser.Game({ type: Phaser.AUTO, parent: host.current, width: GAME_WIDTH, height: GAME_HEIGHT, pixelArt: true, backgroundColor: "#20344d", physics: { default: "arcade", arcade: { debug: false } }, scene: MeadowScene, scale: { mode: Phaser.Scale.FIT, autoCenter: Phaser.Scale.CENTER_BOTH } });
    return () => game.destroy(true);
  }, []);
  return <div className="game-canvas" ref={host} aria-label="Playable meadow tile map" />;
}

export default GameCanvas;
