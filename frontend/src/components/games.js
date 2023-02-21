import {useEffect, useState} from "react";
import axios from "axios";
import Container from "react-bootstrap/Container";
import {Button, Pagination, Table} from "react-bootstrap";

function Games() {
    const [games, setGames] = useState([{user: "Alice", player_1: "Alice", player_2: "Bob", id: 0}]);

    useEffect(() => {
        axios.get("/api/games").then(r => setGames(r.data)).catch(e => console.log(e));
    }, [])

    function joinButtonVariant(game) {
        let variant = "success";
        if (game.player_1 !== game.user && game.player_2 === "") {
            variant = "danger"; // you will join this game
        } else if (game.player_1 !== game.user && game.player_2 !== game.user) {
            variant = "secondary"; // you will only view this game
        }
        return variant;
    }

    function joinButtonText(game) {
        let btnText = "Continue";
        if (game.player_1 !== game.user && game.player_2 === "") {
            btnText = "Join"; // you will join this game
        } else if (game.player_1 !== game.user && game.player_2 !== game.user) {
            btnText = "View"; // you will only view this game
        }
        return btnText;
    }

    return (
        <Container>
            <h1>Open Games List</h1>
            <p>This is a list of all the games that are currently open.</p>
            <p>Click on a game to view it.</p>
            <Table striped bordered hover size="sm" variant="light" responsive className="opacity-75 text-opacity-100">
                <thead className="bg-primary">
                <th>Player1</th>
                <th>Player2</th>
                <th>Action</th>
                </thead>
                <tbody>
                {games.map(game => (
                    <tr key={game.id}>
                        <td>{game.player_1}</td>
                        <td>{game.player_2}</td>
                        <td>
                            <Button variant={joinButtonVariant(game)} href={`/games/${game.id}`}>
                                {joinButtonText(game)}
                            </Button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </Table>
            <Pagination>
                <Pagination.First />
                <Pagination.Prev />
                <Pagination.Next />
                <Pagination.Last />
            </Pagination>
            <Button variant="primary" href={`/games/new`}>New Game</Button>
        </Container>
    );
}

export default Games;