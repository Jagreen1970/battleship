import './App.css';

import 'bootstrap/dist/css/bootstrap.min.css';

import {BrowserRouter, Route, Routes} from "react-router-dom";
import Container from "react-bootstrap/Container";

import {memo} from "react";

import Layout from "./components/layout";
import Games from "./components/games";
import Game from "./components/game";
import Scoreboard from "./components/scoreboard";

function About() {
    return (
        <Container>
            <h1>About</h1>
        </Container>
    )
}

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout/>}>
                    <Route exact path="/games" element={<Games/>}/>
                    <Route path="/games/:id" element={<Game/>}/>
                    <Route path="/scoreboard" element={<Scoreboard/>}/>
                    <Route path="/about" element={<About/>}/>
                </Route>
            </Routes>
        </BrowserRouter>
    );
}

export default memo(App);
