import React from "react";
import { Table, Form, Button } from "react-bootstrap";

const ShipsBoard = ({ boardData, handleShot }) => {
    const renderBoard = () => {
        return boardData.map((row, rowIndex) => {
            return (
                <tr key={rowIndex}>
                    {row.map((cell, cellIndex) => {
                        return (
                            <td key={cellIndex}>
                                {cell === 1 ? (
                                    <Form.Check type="checkbox" disabled checked />
                                ) : (
                                    <Form.Check type="checkbox" disabled />
                                )}
                            </td>
                        );
                    })}
                </tr>
            );
        });
    };

    return (
        <Table striped bordered hover>
            <tbody>{renderBoard()}</tbody>
        </Table>
    );
};

export default ShipsBoard;
