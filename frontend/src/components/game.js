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
            maps: [
                {title: "Alison", map: DummyShipsMapAlison},
                {title: "Bob", map: DummyShotsMapAlison}
            ],
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

    console.log(game_id);

    useEffect(() => {
        axios.get(`/api/games/${game_id}`)
            .then(r => {
                console.log(r.data)
                setGameStats(r.data)
            })
            .catch(e => console.log(e))
    }, [game_id])

    function shipCellClicked(row_index, cell_index) {
        return () => {
            if(gameStats.status === "1") {
                return
            }
            let call;
            if(gameStats.board.maps["Ships"][row_index][cell_index] === " ") {
                call = axios.put("/api/games/" + gameStats._id + "/pin/" + row_index + "-" + cell_index)
            } else if(gameStats.board.ships_map[row_index][cell_index] === "O") {
                call = axios.delete("/api/games/" + gameStats._id + "/pin/" + row_index + "-" + cell_index)
            }
            call.then((res) => {
                console.log(res)
                if(res.status === 200) {
                    setGameStats(res.data)
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

            if(gameStats.board.maps[1].map[row_index][cell_index] !== " ") {
                // already shot on that one
                return;
            }

            axios.post("/api/games/" + gameStats._id + "/target", {
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

    function gameStateSetup(status) {
        return status === "0"
    }

    function gameStatePlying(status) {
        return status === "1"
    }

    function isMyTurn(user) {
        return gameStats.player_to_move === user
    }

    function isPlayer(user) {
        return user === gameStats.player_1.name || user === gameStats.player_2.name;
    }

    function isOpponentMap(map_name) {
        if (!isPlayer(gameStats.user)) {
            return false
        }
        return map_name !== gameStats.user
    }

    function onCellClicked (map_name) {
        return (row_index, cell_index) => {
            if (gameStatePlying(gameStats.status)
                && isPlayer(gameStats.user)
                && isMyTurn(gameStats.user)
                && isOpponentMap(map_name)) {
                return targetCellClicked(row_index, cell_index)
            }

            if (gameStats.status === "0" && map_name === gameStats.user) {
                return shipCellClicked(row_index, cell_index)
            }
        }
    }

    function isMapDisabled (map_name) {
        return () => {
            if (!isPlayer(gameStats.user)) {
                return true
            }

            if (gameStatePlying(gameStats.status) && isOpponentMap(map_name)) {
                return false
            }

            return !(gameStateSetup(gameStats.status) && !isOpponentMap(map_name));

        }
    }

    return (
        <Container>
            <GameStats {...gameStats}/>
            <GameBoard
                isStatusPlaying={gameStats.status === "1"}
                my_turn={gameStats.player_to_move === gameStats.user}
                onCellClicked={onCellClicked}
                isMapDisabled={isMapDisabled}
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