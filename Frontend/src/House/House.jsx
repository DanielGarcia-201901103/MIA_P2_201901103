import './House.css';
import Navb from "../Nav/Navb.jsx"
const House = () => {
  return (
    <>
      <Navb />
      <h1>¡Hola! &#x1F44B;</h1>
      <h2>Bienvenido al Proyecto </h2>
      <h3>Esperamos tengas una buena experiencia</h3>
      <div className="container">
        <div className="loader"></div>
        <div className="loader"></div>
        <div className="loader"></div>
      </div>

      <footer>
        Josué Daniel Rojché García <br />
        201901103
      </footer>
    </>
  );
}

export default House;
