import Cell from "./cell";

const Board = (props) => {
    let letter = "ABCDEFGHIJ";
    return (
        <>
            <h2 className={"game-board-title"}>{props.boardTitle}</h2>
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
                    {props.ships_map.map((row, row_index) => (
                        <div className="map-row" key={row_index}>
                            <div className="map-index-cell">{row_index + 1}</div>
                            {row.map((cell, cell_index) => (
                                <Cell value={cell} disabled={props.isMapDisabled()}
                                      onClick={props.onCellClicked(row_index, cell_index)}/>
                            ))}
                        </div>
                    ))
                    }
                </div>
            </div>
        </>
    )
}


export default Board
