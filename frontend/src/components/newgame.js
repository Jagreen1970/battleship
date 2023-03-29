import Container from "react-bootstrap/Container";
import {Button, Form, Modal} from "react-bootstrap";
import {useState} from "react";
import UserSession from "./clientsession";
import {useNavigate} from "react-router-dom";
import axios from "axios";

function NewGameModal() {
    const [showNewGame, setShowNewGame] = useState(false);
    const [playerName, setPlayerName] = useState("");

    let navigate = useNavigate();

    const handleCloseNewGame = () => setShowNewGame(false);
    const handleShowNewGame = () => {
        if (UserSession.getUserName() === "") {
            setShowNewGame(true);
            return;
        }

        axios.post("/api/games", {})
            .then(res => {
                console.log(res);
                navigate(`/game/${res.data._id}`);
            })
            .catch(e => console.log(e));
    }
    const handleCreateNewGame = () => {
        setShowNewGame(false);
        UserSession.setUserName(playerName);
        axios.post("/api/login", {username: playerName})
            .then(() => {
                return axios.post("/api/games", {})
            })
            .then(res => {
                console.log(res);
                navigate(`/game/${res.data._id}`);
            })
            .catch(e => console.log(e));
    }

    return (
        <Container>
            <Button variant="primary" onClick={handleShowNewGame}>New Game</Button>
            <Modal
                show={showNewGame}
                onHide={handleCloseNewGame}
                size="lg"
                aria-labelledby="contained-modal-title-vcenter"
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title id="contained-modal-title-vcenter">Enter your name first!</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form>
                        <Form.Group>
                            <Form.Label>Name</Form.Label>
                            <Form.Control type="text" onChange={(e) => setPlayerName(e.target.value)}/>
                        </Form.Group>
                    </Form>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={handleCloseNewGame}>Close</Button>
                    <Button variant="primary" onClick={handleCreateNewGame}>Create</Button>
                </Modal.Footer>
            </Modal>
        </Container>
    )
}

export default NewGameModal;