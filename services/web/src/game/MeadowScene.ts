import * as Phaser from "phaser";
import { BOULDER_POSITIONS, GAME_HEIGHT, GAME_WIDTH, MAP_HEIGHT, MAP_WIDTH, PLAYER_SPEED, TILE_SIZE } from "../lib/game-constants";

type MovementKeys = { W: Phaser.Input.Keyboard.Key; A: Phaser.Input.Keyboard.Key; S: Phaser.Input.Keyboard.Key; D: Phaser.Input.Keyboard.Key; };

class MeadowScene extends Phaser.Scene {
  private player?: Phaser.Physics.Arcade.Sprite;
  private cursors?: Phaser.Types.Input.Keyboard.CursorKeys;
  private keys?: MovementKeys;

  constructor() { super("meadow"); }

  create() {
    this.createTileTexture();
    this.createPlayer();
    this.createMap();
    this.add.text(16, 16, "MEADOW TRAIL", { fontFamily: "monospace", fontSize: "16px", color: "#ffffff", stroke: "#20344d", strokeThickness: 4 }).setDepth(3);
  }

  update() {
    if (!this.player || !this.cursors || !this.keys) return;
    const horizontal = Number(this.cursors.right.isDown || this.keys.D.isDown) - Number(this.cursors.left.isDown || this.keys.A.isDown);
    const vertical = Number(this.cursors.down.isDown || this.keys.S.isDown) - Number(this.cursors.up.isDown || this.keys.W.isDown);
    const movement = new Phaser.Math.Vector2(horizontal, vertical).normalize().scale(PLAYER_SPEED);
    this.player.setVelocity(movement.x, movement.y);
  }

  private createTileTexture() {
    const tiles = this.add.graphics();
    tiles.fillStyle(0x82b949).fillRect(0, 0, TILE_SIZE, TILE_SIZE);
    tiles.fillStyle(0x70a942).fillRect(0, 1, TILE_SIZE, 2);
    tiles.fillStyle(0x689d3b).fillRect(0, TILE_SIZE - 2, TILE_SIZE, 2);
    tiles.fillStyle(0x3b628d).fillRect(TILE_SIZE, 0, TILE_SIZE, TILE_SIZE);
    tiles.fillStyle(0x4778a8).fillRect(TILE_SIZE, 2, TILE_SIZE, 3);
    tiles.fillStyle(0x604d39).fillRect(TILE_SIZE * 2, 0, TILE_SIZE, TILE_SIZE);
    tiles.fillStyle(0x7b654a).fillRect(TILE_SIZE * 2 + 3, 3, TILE_SIZE - 6, TILE_SIZE - 6);
    tiles.generateTexture("tiles", TILE_SIZE * 3, TILE_SIZE);
    tiles.destroy();
  }

  private createMap() {
    const map = this.make.tilemap({ width: MAP_WIDTH, height: MAP_HEIGHT, tileWidth: TILE_SIZE, tileHeight: TILE_SIZE });
    const tileSet = map.addTilesetImage("tiles", "tiles", TILE_SIZE, TILE_SIZE, 0, 0)!;
    const ground = map.createBlankLayer("ground", tileSet, 0, 0)!;
    const water = map.createBlankLayer("water", tileSet, 0, 0)!;
    const boulders = map.createBlankLayer("boulders", tileSet, 0, 0)!;
    for (let y = 0; y < MAP_HEIGHT; y += 1) for (let x = 0; x < MAP_WIDTH; x += 1) ground.putTileAt(0, x, y);
    for (let x = 0; x < MAP_WIDTH; x += 1) { water.putTileAt(1, x, 0); water.putTileAt(1, x, MAP_HEIGHT - 1); }
    for (let y = 1; y < MAP_HEIGHT - 1; y += 1) { water.putTileAt(1, 0, y); water.putTileAt(1, MAP_WIDTH - 1, y); }
    BOULDER_POSITIONS.forEach(([x, y]) => boulders.putTileAt(2, x, y));
    water.setCollision(1); boulders.setCollision(2);
    this.createPlayerCollisions(water, boulders);
  }

  private createPlayer() {
    const hero = this.add.graphics();
    hero.fillStyle(0x26314d).fillRect(4, 11, 24, 17);
    hero.fillStyle(0xf0b389).fillCircle(16, 10, 9);
    hero.fillStyle(0xe75849).fillRect(6, 1, 20, 8);
    hero.fillStyle(0x172238).fillRect(7, 28, 7, 4).fillRect(18, 28, 7, 4);
    hero.generateTexture("hero", TILE_SIZE, TILE_SIZE); hero.destroy();
    this.player = this.physics.add.sprite(TILE_SIZE * 3.5, TILE_SIZE * 3.5, "hero");
    this.player.setCollideWorldBounds(true).setDepth(2);
    this.physics.world.setBounds(0, 0, GAME_WIDTH, GAME_HEIGHT);
    this.cursors = this.input.keyboard?.createCursorKeys();
    this.keys = this.input.keyboard?.addKeys("W,A,S,D") as MovementKeys;
  }

  private createPlayerCollisions(water: Phaser.Tilemaps.TilemapLayer, boulders: Phaser.Tilemaps.TilemapLayer) {
    if (!this.player) return;
    this.physics.add.collider(this.player, water);
    this.physics.add.collider(this.player, boulders);
  }
}

export default MeadowScene;
