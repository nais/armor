import React from "react";

const style = {
    backgroundColor: "#F8F8F8",
    borderTop: "1px solid #E7E7E7",
    textAlign: "center",
    padding: "20px",
    position: "fixed",
    left: "0",
    bottom: "0",
    height: "40px",
    width: "100%"
};

function Footer() {
    return (
        <a style={style} className='text-dark' href='https://naas.nais.io/'>
            naas.nais.io
        </a>
    );
}

export default Footer;