import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";

import Spinner from "../Spinner/Spinner";

export const Overlay = styled.div`
  width: 100%;
  height: 100%;
  z-index: 1;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  position: absolute;
  background: rgba(247, 247, 247, 0.95);

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
  align-items: center;

  opacity: ${props => (props.show ? 1 : 0)};
  visibility: ${props => (props.show ? "visible" : "hidden")};
  transition: all 0.25s ease-in-out;

  > div {
    transform: scale(1.3) translateY(0rem);
  }
`;

export default function PageLoader({ show }) {
  return (
    <Overlay show={show}>
      <Spinner />
    </Overlay>
  );
}

PageLoader.propTypes = {
  /* Show or hide loader  */
  show: PropTypes.bool
};

PageLoader.defaultProps = {
  show: false
};
