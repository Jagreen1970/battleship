import {Col, Container, Row} from "react-bootstrap";
import "./gameboard.css";
import hit from "../images/Hit.png";
import miss from "../images/Miss.png";
import pin from "../images/Pin.png";

const Cell = (props) => {
    return (
        <div className="map-cell"
             onClick={props.onClick}>
            {props.value === "X" && <img src={hit} alt="hit" className="pin"/>}
            {props.value === "-" && <img src={miss} alt="miss" className="pin"/>}
            {props.value === "O" && <img src={pin} alt="pin" className="pin"/>}
        </div>
    )
}

const GameBoard = (board) => {
    let letter = "ABCDEFGHIJ";

    return (
        <Container fluid>
            <Row>
                <Col sm={6}>
                    <h2>Ships</h2>
                    <div className="ships-map">
                        <div className="map-index-row">
                            <div className="map-row">
                                <div className="map-index-cell">#</div>
                                {letter.split("").map((letter, index) => (
                                    <div className="map-index-cell">{letter}</div>
                                ))}
                            </div>
                        </div>
                        <div className="map-body">
                            {board.ships_map.map((ship, row_index) => (
                                <div className="map-row" key={row_index}>
                                    <div className="map-index-cell">{row_index + 1}</div>
                                    {board.ships_map[row_index].map((cell, cell_index) => (
                                        <Cell value={cell} disabled={board.isStatusPlaying}
                                              onClick={board.shipCellClicked(row_index, cell_index)}/>
                                    ))}
                                </div>
                            ))
                            }
                        </div>
                    </div>
                </Col>
                <Col sm={6}>
                    <h2>Target</h2>
                    <div className="ships-map">
                        <div className="map-index-row">
                            <div className="map-row">
                                <div className="map-index-cell">#</div>
                                {letter.split("").map((letter, index) => (
                                    <div className="map-index-cell">{letter}</div>
                                ))}
                            </div>
                        </div>
                        <div className="map-body">
                            {board.shots_map.map((target, row_index) => (
                                <div className="map-row" key={row_index}>
                                    <div className="map-index-cell">{row_index + 1}</div>
                                    {board.shots_map[row_index].map((cell, cell_index) => (
                                        <Cell value={cell} disabled={!board.isStatusPlaying || !board.my_turn}
                                              onClick={board.targetCellClicked(row_index, cell_index)}/>
                                    ))}
                                </div>
                            ))
                            }
                        </div>
                    </div>
                </Col>
            </Row>
        </Container>
    )
}

export default GameBoard