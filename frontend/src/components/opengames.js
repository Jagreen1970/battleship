import Container from "react-bootstrap/Container";
import {Button, Pagination, Table} from "react-bootstrap";
import {useEffect, useState} from "react";
import axios from "axios";
import {Link} from "react-router-dom";

export function OpenGamesList(props) {

    const [gamesList, setGamesList] = useState({
        user: "Alice",
        games: [{
            player_1: {
                name: "Alice",
            },
            player_2: {
                name: "Bob",
            },
            _id: 0
        }]
    });

    useEffect(() => {
        axios.get("/api/games")
            .then(r => {
                console.log(r.data);
                setGamesList(r.data)
            })
            .catch(e => console.log(e));
    }, [])

    function isPlayer(game, user) {
        return game.player_1.name === user || (game.player_2 !== null && game.player_2.name === user);
    }

    function isGuest(user) {
        return user === null || user === "guest"
    }

    function joinButtonVariant(game, user) {
        if (isPlayer(game, user)) {
            return "success"
        }
        if (isGuest(user)) {
            return "info"
        }
        return "primary"
    }

    function joinButtonText(game, user) {
        if (isPlayer(game, user)) {
            return "Continue"
        }
        if (isGuest(user)) {
            return "View"
        }
        return "Join"
    }

    return (
        <Container>
            <h1>Open Games List</h1>
            <p>This is a list of all the games that are currently open.</p>
            <p>Click the respective button to view or join a game.</p>
            <Table striped bordered hover size="sm" variant="light" responsive className="opacity-75 text-opacity-100">
                <thead className="bg-primary">
                <tr>
                    <th>Player1</th>
                    <th>Player2</th>
                    <th>Action</th>
                </tr>
                </thead>
                <tbody>
                {gamesList.games.map(game => (
                    <tr key={game._id}>
                        <td>{game.player_1.name}</td>
                        <td>{"" || (game.player_2 && game.player_2.name)}</td>
                        <td>
                            <Button variant={joinButtonVariant(game, gamesList.user)} as={Link} to={`/games/${game._id}`}>
                                {joinButtonText(game, gamesList.user)}
                            </Button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </Table>
            <Pagination>
                <Pagination.First/>
                <Pagination.Prev/>
                <Pagination.Next/>
                <Pagination.Last/>
            </Pagination>
        </Container>
    );
}
