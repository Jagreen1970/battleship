import GameBoard from "./gameboard";
import Container from "react-bootstrap/Container";
import {memo, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import axios from "axios";


function GameStats(game) {
    return (
        <Container>
            <h2 className="text-center">Game Stats</h2>
            <h3 className="text-center">{game.player_1.name} vs. {game.player_2.name}</h3>
        </Container>
    )
}

const Game = () => {
    const {game_id} = useParams()
    const [gameStats, setGameStats] = useState({
        _id: 0,
        user: "Alison",
        player_1: {
            id: "",
            name: "Alison",
            score: 0,
        },
        player_2: {
            id: "",
            name: "Bob",
            score: 0,
        },
        player_to_move: "Alison",
        history: [{player: "Alison", hit: true, x: 1, y: 1}, {player: "Bob", hit: false, x: 0, y: 0}],
        status: "1",
        board: {
            pins_available: 0,
            shots_map: DummyShotsMapAlison,
            ships_map: DummyShipsMapAlison,
            fleet: [
                {
                    ship_type: "Submarine",
                    fields: [{x: 1, y: 1}, {x: 1, y: 2}],
                    orientation: "Vertical"
                },
                {
                    ship_type: "Submarine",
                    fields: [{x: 3, y: 1}, {x: 4, y: 1}],
                    orientation: "Horizontal"
                }
            ]
        },
    });

    function fetchStats() {
        fetch(`/api/games/${game_id}`)
            .then(res => res.json())
            .then(data => {
                console.log(data)
                setGameStats(data)
            })
    }

    useEffect(() => {
        fetchStats()
    }, [])

    function shipCellClicked(row_index, cell_index) {
        return () => {
            if(gameStats.status === "1") {
                return
            }
            let call;
            if(gameStats.board.ships_map[row_index][cell_index] === " ") {
                call = axios.put("/api/game/" + gameStats._id + "/pin/" + row_index + "-" + cell_index)
            } else if(gameStats.board.ships_map[row_index][cell_index] === "O") {
                call = axios.delete("/api/game/" + gameStats._id + "/pin/" + row_index + "-" + cell_index)
            }
            call.then((res) => {
                if(res.status === 200) {
                    setGameStats(res.data)
                } else {
                    console.log(res)
                }
            }).catch(
                (error) => {
                    console.log(error)
                }
            )
        }
    }

    function targetCellClicked(row_index, cell_index) {
        return () => {
            if(gameStats.status !== "1" || gameStats.player_to_move !== gameStats.user) {
                return
            }

            if(gameStats.board.ships_map[row_index][cell_index] !== " ") {
                return;
            }

            axios.post("/api/game/" + gameStats._id + "/target", {
                x: row_index,
                y: cell_index
            }).then((res) => {
                if(res.status === 200) {
                    setGameStats(res.data)
                } else {
                    console.log(res)
                }
            }).catch(
                (error) => {
                    console.log(error)
                }
            )
        }
    }

    return (
        <Container>
            <GameStats {...gameStats}/>
            <GameBoard
                isStatusPlaying={gameStats.status === "1"}
                my_turn={gameStats.player_to_move === gameStats.user}
                shipCellClicked={shipCellClicked}
                targetCellClicked={targetCellClicked}
                {...gameStats.board}/>
        </Container>
    )
}

let DummyShipsMapAlison = [
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", "O", " ", "O", "O", " ", " ", " ", " ", " "],
    [" ", "O", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
];

let DummyShotsMapAlison = [
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", "X", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
    [" ", " ", " ", " ", " ", " ", " ", " ", " ", " "],
];


export default memo(Game);