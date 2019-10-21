import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";

const StyledSpinner = styled.div`
  position: relative;
  transform: translateY(${props => props.translateY}rem);
  > div {
    background-color: ${props => props.color};
    border-radius: 100%;
    margin: 2px;
    animation-fill-mode: both;
    position: absolute;
    left: -${props => props.size / 2}px;
    top: 0;
    opacity: 1;
    margin: 0;
    width: ${props => props.size}px;
    height: ${props => props.size}px;
    animation: ball-scale-multiple 1s 0s linear infinite;
  }

  > div:nth-child(2) {
    animation-delay: -0.4s;
  }

  > div:nth-child(3) {
    animation-delay: -0.2s;
  }

  @keyframes ball-scale-multiple {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    5% {
      opacity: 1;
    }
    100% {
      transform: scale(1);
      opacity: 0;
    }
  }
`;

const Spinner = ({ size }) => {
  return (
    <StyledSpinner size={size} color="#333">
      <div />
      <div />
      <div />
    </StyledSpinner>
  );
};

Spinner.propTypes = {
  /**
   * Size of button to be rendered, e.g. 10,20,..., 50, 100 etc.
   */
  size: PropTypes.number
};

Spinner.defaultProps = {
  size: 60
};

export default Spinner;
