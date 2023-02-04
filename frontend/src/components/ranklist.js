import React, { useState, useEffect } from "react";
import axios from "axios";
import { Table, Pagination } from "react-bootstrap";

const PlayerRankList = () => {
    const [players, setPlayers] = useState([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [playersPerPage, setPlayersPerPage] = useState(10);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get("/api/players");
                setPlayers(response.data.slice(0, 100));
            } catch (error) {
                console.error(error);
            }
        };
        fetchData();
    }, []);

    // Berechne den ersten und den letzten Spieler auf der aktuellen Seite
    const indexOfLastPlayer = currentPage * playersPerPage;
    const indexOfFirstPlayer = indexOfLastPlayer - playersPerPage;
    const currentPlayers = players.slice(indexOfFirstPlayer, indexOfLastPlayer);

    // Berechne die Anzahl der Seiten basierend auf der Anzahl der Spieler und der Anzahl der Spieler pro Seite
    const pageNumbers = [];
    for (let i = 1; i <= Math.ceil(players.length / playersPerPage); i++) {
        pageNumbers.push(i);
    }

    return (
        <>
            <Table striped bordered hover>
                <thead>
                <tr>
                    <th>Rank</th>
                    <th>Name</th>
                    <th>Score</th>
                </tr>
                </thead>
                <tbody>
                {currentPlayers.map((player, index) => (
                    <tr key={player.id}>
                        <td>{indexOfFirstPlayer + index + 1}</td>
                        <td>{player.name}</td>
                        <td>{player.score}</td>
                    </tr>
                ))}
                </tbody>
            </Table>
            <div className="d-flex justify-content-center mt-3">
                <Pagination>
                    <Pagination.First onClick={() => setCurrentPage(1)} />
                    <Pagination.Prev
                        onClick={() => setCurrentPage(currentPage - 1)}
                        disabled={currentPage === 1}
                    />
                    {pageNumbers.map((number) => (
                        <Pagination.Item
                            key={number}
                            active={number === currentPage}
                            onClick={() => setCurrentPage(number)}
                        >
                            {number}
                        </Pagination.Item>
                    ))}
                    <Pagination.Next
                        onClick={() => setCurrentPage(currentPage + 1)}
                        disabled={currentPage === pageNumbers[pageNumbers.length - 1]}
                    />
                    <Pagination.Last
                        onClick={() => setCurrentPage(pageNumbers[pageNumbers.length - 1])}
                    />
                </Pagination>
            </div>
            <div className="d-flex justify-content-center mt-3">
                <Pagination.Item
                    onClick={() => setPlayersPerPage(10)}
                    active={playersPerPage === 10}
                >
                    10
                </Pagination.Item>
                <Pagination.Item
                    onClick={() => setPlayersPerPage(25)}
                    active={playersPerPage === 25}
                >
                    25
                </Pagination.Item>
                <Pagination.Item
                    onClick={() => setPlayersPerPage(50)}
                    active={playersPerPage === 50}
                >
                    50
                </Pagination.Item>
            </div>
        </>
    );
};

export default PlayerRankList;
