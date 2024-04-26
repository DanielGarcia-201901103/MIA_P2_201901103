import React from "react";
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import Navb from "./Nav/Navb.jsx"
import House from "./House/House.jsx"
import Pantalla1 from "./Pantalla1/Pantalla1.jsx"
import Pantalla2 from "./Pantalla2/Pantalla2.jsx"
import Pantalla3 from "./Pantalla3/Pantalla3.jsx"

const App = () => {

  return (
    <Router>
      <Routes>
        <Route exact path='/home' element={<House/>} />
        <Route exact path='/pantalla1' element={<Pantalla1/>} />
        <Route exact path='/pantalla2' element={<Pantalla2/>} />
        <Route exact path='/pantalla3' element={<Pantalla3/>} />
      </Routes>
    </Router>
  );
};
export default App;