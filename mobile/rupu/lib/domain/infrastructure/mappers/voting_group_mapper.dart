import 'package:rupu/domain/entities/voting_group_action_result.dart';

class VotingGroupMapper {
  static VotingGroupActionResult fromJson(Map<String, dynamic> json) {
    return VotingGroupActionResult(
      success: json['success'] ?? false,
      message: json['message'],
    );
  }
}
