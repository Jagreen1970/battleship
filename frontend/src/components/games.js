import Container from "react-bootstrap/Container";
import {Row} from "react-bootstrap";
import {OpenGamesList} from "./opengames";
import NewGameModal from "./newgame";

function Games() {
    return (
        <Container>
            <Row>
                <OpenGamesList/>
            </Row>
            <Row>
                <NewGameModal/>
            </Row>
        </Container>);
}

export default Games;