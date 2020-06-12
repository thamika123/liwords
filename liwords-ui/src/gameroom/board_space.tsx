import React from 'react';
import { ReactComponent as Logo } from '../assets/aero.svg';
import { BonusType } from '../constants/board_layout';
import { ArrowRightOutlined, ArrowDownOutlined } from '@ant-design/icons/lib';

interface BonusProperties {
  fillColor: string;
  bonusText: string;
}

function getBonusProperties(bt: BonusType): BonusProperties {
  switch (bt) {
    case BonusType.DoubleWord:
      return { fillColor: '#FFC9F3', bonusText: '2WS' };
    case BonusType.TripleWord:
      return { fillColor: '#FD7070', bonusText: '3WS' };
    case BonusType.DoubleLetter:
      return { fillColor: '#C9E5FF', bonusText: '2LS' };
    case BonusType.TripleLetter:
      return { fillColor: '#6F87DF', bonusText: '3LS' };
  }
  return { fillColor: 'hsl(35, 30%, 98%)', bonusText: '' };
}

type Props = {
  bonusType: BonusType;
  boardSquareDim: number;
  x: number;
  y: number;
  showBonusLabel: boolean;
  startingSquare: boolean;
  arrow: boolean;
  arrowHoriz: boolean;
  clicked: () => void;
};

const BoardSpace = (props: Props) => {
  const { fillColor, bonusText } = getBonusProperties(props.bonusType);

  let bonusLabel = null;
  let startingSquare = null;
  let arrow = null;
  if (props.showBonusLabel && bonusText !== '') {
    bonusLabel = <p className="bonus-label">{bonusText}</p>;
  }
  // ✩✪✫
  if (props.startingSquare) {
    startingSquare = (
      <Logo
        className="logo"
        width={props.boardSquareDim / 1.5}
        height={props.boardSquareDim / 1.5}
      />
    );
  }
  if (props.arrow) {
    if (props.arrowHoriz) {
      arrow = <ArrowRightOutlined />;
    } else {
      arrow = <ArrowDownOutlined />;
    }
  }

  const styleOverrides = {
    backgroundColor: fillColor,
  };

  return (
    <div className="board-space" onClick={props.clicked} style={styleOverrides}>
      {bonusLabel}
      {startingSquare}
      {arrow}
    </div>
  );
};

export default BoardSpace;