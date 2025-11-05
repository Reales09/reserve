class Voting {
  final int id;
  final int votingGroupId;
  final String title;
  final String description;
  final String votingType;
  final bool isSecret;
  final bool allowAbstention;
  final bool isActive;
  final int displayOrder;
  final int requiredPercentage;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Voting({
    required this.id,
    required this.votingGroupId,
    required this.title,
    required this.description,
    required this.votingType,
    required this.isSecret,
    required this.allowAbstention,
    required this.isActive,
    required this.displayOrder,
    required this.requiredPercentage,
    required this.createdAt,
    required this.updatedAt,
  });
}
