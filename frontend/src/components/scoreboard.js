import React, {useState, useEffect} from "react";
import axios from "axios";
import {Container, Table} from "react-bootstrap";

const PlayerRankList = () => {
    const [players, setPlayers] = useState({
        scores: [
            {
                id: 1,
                name: "John Doe",
                score: 100
            }
        ]
    });

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get("/api/scoreboard");
                console.log(response.data);
                setPlayers(response.data);
            } catch (error) {
                console.error(error);
            }
        };
        fetchData();
    }, []);

    return (
        <Container>
            <Table striped bordered hover size="sm" variant="light" responsive className="opacity-75 text-opacity-100">
                <thead className="bg-primary">
                <tr>
                    <th>Rank</th>
                    <th>Name</th>
                    <th>Score</th>
                </tr>
                </thead>
                <tbody>
                {players.scores.map((player, index) => (
                    <tr key={player.id}>
                        <td>{index + 1}</td>
                        <td>{player.name}</td>
                        <td>{player.score}</td>
                    </tr>
                ))}
                </tbody>
            </Table>
        </Container>
    );
};

export default PlayerRankList;
