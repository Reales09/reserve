import 'package:rupu/domain/entities/voting.dart';
import 'package:rupu/domain/entities/votings_result.dart';

class VotingMapper {
  static VotingsResult fromJson(Map<String, dynamic> json) {
    return VotingsResult(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<Voting>.from(
        (json['data'] as List<dynamic>).map(
          (x) => Voting(
            id: x['id'],
            votingGroupId: x['voting_group_id'],
            title: x['title'],
            description: x['description'],
            votingType: x['voting_type'],
            isSecret: x['is_secret'],
            allowAbstention: x['allow_abstention'],
            isActive: x['is_active'],
            displayOrder: x['display_order'],
            requiredPercentage: x['required_percentage'],
            createdAt: DateTime.parse(x['created_at']),
            updatedAt: DateTime.parse(x['updated_at']),
          ),
        ),
      ),
    );
  }
}
