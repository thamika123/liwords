import { parseMsgs } from './socket_handlers';
import {
  MessageType,
  ServerGameplayEvent,
  ServerChallengeResultEvent,
} from '../gen/api/proto/realtime/realtime_pb';
import {
  GameEvent,
  ChallengeRule,
} from '../gen/macondo/api/proto/macondo/macondo_pb';

const toArr = (s: string) => {
  const bytes = new Uint8Array(Math.ceil(s.length / 2));
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(s.substr(i * 2, 2), 16);
  }
  return bytes;
};

it('parses multiple messages', () => {
  const multimsg =
    '003d040a1e0a0962726f776e2d6d656e200152085245522e46554e4478629001' +
    'fef91b1216656b584a61326d5676587a4c43477535667448457a4c20fef91b00' +
    '1b091216786a4357756737455a7444784448583566525a544c6f1803';
  const msg = toArr(multimsg);
  const res = parseMsgs(msg);
  expect(res.length).toBe(2);
  expect(res[0].msgType).toBe(MessageType.SERVER_GAMEPLAY_EVENT);
  expect(res[1].msgType).toBe(MessageType.SERVER_CHALLENGE_RESULT_EVENT);

  const firstEvt = new ServerGameplayEvent();
  console.log(res[0].parsedMsg);
  firstEvt.setTimeRemaining(457982);
  firstEvt.setGameId('ekXJa2mVvXzLCGu5ftHEzL');
  const inevt = new GameEvent();
  inevt.setNickname('brown-men'); // this was an autogenerated name, don't @ me, cancel culture
  inevt.setType(GameEvent.Type.PHONY_TILES_RETURNED);
  /**
   * 			Nickname:    lastEvent.Nickname,
			Type:        pb.GameEvent_PHONY_TILES_RETURNED,
			LostScore:   lastEvent.Score,
			Cumulative:  cumeScoreBeforeChallenge - lastEvent.Score,
			Rack:        lastEvent.Rack,
			PlayedTiles: lastEvent.PlayedTiles,
			// Note: these millis remaining would be the challenger's
			MillisRemaining: int32(millis),
   */
  inevt.setLostScore(98);
  inevt.setPlayedTiles('RER.FUND');
  inevt.setMillisRemaining(457982);
  firstEvt.setEvent(inevt);
  expect(res[0].parsedMsg).toEqual(firstEvt);

  const secondEvt = new ServerChallengeResultEvent();
  secondEvt.setChallenger('xjCWug7EZtDxDHX5fRZTLo');
  secondEvt.setChallengeRule(ChallengeRule.FIVE_POINT);
  expect(res[1].parsedMsg).toEqual(secondEvt);
  // Even though it is not specified as false, it should still be false.
  expect((res[1].parsedMsg as ServerChallengeResultEvent).getValid()).toBe(
    false
  );
});