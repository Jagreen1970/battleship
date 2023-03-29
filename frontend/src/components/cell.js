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

export default Cell;