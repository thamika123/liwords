import React from 'react';
import { Button } from 'antd';
import { PlayerAvatar } from '../shared/player_avatar';
import { PlayerMetadata } from '../gameroom/game_info';

type Props = {
  player: Partial<PlayerMetadata> | undefined;
  handleLogout?: () => void;
};

export const LogOut = React.memo((props: Props) => {
  return (
    <div className="log-out">
      <h3>Log out</h3>
      <div className="avatar-container">
        <PlayerAvatar player={props.player} />
        <div className="full-name">{props.player?.full_name}</div>
      </div>
      <div>
        You’ll have to log back in to your account to play games or see tiles
        while watching tournament games on Woogles.io.
      </div>
      <Button type="primary" onClick={props.handleLogout}>
        Log out
      </Button>
    </div>
  );
});
