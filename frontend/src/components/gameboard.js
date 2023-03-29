import {Col, Container, Row} from "react-bootstrap";
import "./gameboard.css";
import Board from "./board";

const GameBoard = (props) => {
    return (
        <Container fluid>
            <Row>
                {props.maps.map(
                    (entry, index) => (
                        <Col key={index} xs={12} md={6}>
                            <Board boardTitle={entry.title}
                                   ships_map={entry.map}
                                   onCellClicked={props.onCellClicked(entry.title)}
                                   isMapDisabled={props.isMapDisabled(entry.title)}/>
                        </Col>
                    )
                )}
            </Row>
        </Container>
    )
}

export default GameBoard