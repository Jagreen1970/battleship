import React, { useState, useEffect } from 'react';
import axios from 'axios';

const ShotBoard = ({ boardData, onShotTaken }) => {
    const [board, setBoard] = useState(boardData);

    const handleShotTaken = async (rowIndex, columnIndex) => {
        const result = await axios.post('/api/fire-shot', {
            row: rowIndex,
            column: columnIndex,
        });

        if (result.data.hit) {
            setBoard(prevBoard =>
                prevBoard.map((row, i) =>
                    i === rowIndex
                        ? row.map((cell, j) => (j === columnIndex ? 'X' : cell))
                        : row
                )
            );
        } else {
            setBoard(prevBoard =>
                prevBoard.map((row, i) =>
                    i === rowIndex
                        ? row.map((cell, j) => (j === columnIndex ? 'M' : cell))
                        : row
                )
            );
        }
    };

    return (
        <table>
            <tbody>
            {board.map((row, rowIndex) => (
                <tr key={rowIndex}>
                    {row.map((cell, cellIndex) => (
                        <td key={cellIndex}>
                            <button onClick={() => handleShotTaken(rowIndex, cellIndex)}>
                                {cell}
                            </button>
                        </td>
                    ))}
                </tr>
            ))}
            </tbody>
        </table>
    );
};

export default ShotBoard;
