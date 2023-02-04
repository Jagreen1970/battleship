import React, { useState, useEffect } from "react";
import { Container, Row, Col, Spinner } from "react-bootstrap";
import axios from "axios";
import ShipsBoard from "./shipsboard";
import ShotBoard from "./shotboard";

const GameBoard = (gameID) => {
    const [boardData, setBoardData] = useState(null);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);
            const response = await axios.get("/api/game");
            setBoardData(response.data);
            setIsLoading(false);
        };
        fetchData();

        const intervalId = setInterval(() => {
            fetchData();
        }, 1000);

        return () => clearInterval(intervalId);
    }, []);

    return (
        <Container>
            {isLoading ? (
                <Spinner animation="border" role="status">
                    <span className="sr-only">Loading...</span>
                </Spinner>
            ) : (
                <Row>
                    <Col xs={6}>
                        <ShipsBoard boardData={boardData} />
                    </Col>
                    <Col xs={6}>
                        <ShotBoard boardData={boardData} />
                    </Col>
                </Row>
            )}
        </Container>
    );
};

export default GameBoard;
