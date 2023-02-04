import React, { useState, useEffect } from "react";
import axios from "axios";
import { Table } from "react-bootstrap";

const PlayerRankList = () => {
    const [players, setPlayers] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get("/api/players");
                setPlayers(response.data);
            } catch (error) {
                console.error(error);
            }
        };
        fetchData();
    }, []);

    return (
        <Table striped bordered hover>
            <thead>
            <tr>
                <th>Rank</th>
                <th>Name</th>
                <th>Score</th>
            </tr>
            </thead>
            <tbody>
            {players.map((player, index) => (
                <tr key={player.id}>
                    <td>{index + 1}</td>
                    <td>{player.name}</td>
                    <td>{player.score}</td>
                </tr>
            ))}
            </tbody>
        </Table>
    );
};

export default PlayerRankList;
